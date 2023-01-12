package backup

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceFramework() *schema.Resource {
	return &schema.Resource{
		Create: resourceFrameworkCreate,
		Read:   resourceFrameworkRead,
		Update: resourceFrameworkUpdate,
		Delete: resourceFrameworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_parameter": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"scope": {
							// The control scope can include
							// one or more resource types,
							// a combination of a tag key and value,
							// or a combination of one resource type and one resource ID.
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compliance_resource_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										MinItems: 1,
										MaxItems: 100,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"compliance_resource_types": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									// A maximum of one key-value pair can be provided.
									// The tag value is optional, but it cannot be an empty string
									"tags": tftags.TagsSchema(),
								},
							},
						},
					},
				},
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validFrameworkName,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},
		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceFrameworkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("name").(string)

	input := &backup.CreateFrameworkInput{
		IdempotencyToken:  aws.String(resource.UniqueId()),
		FrameworkControls: expandFrameworkControls(d.Get("control").(*schema.Set).List()),
		FrameworkName:     aws.String(name),
	}

	if v, ok := d.GetOk("description"); ok {
		input.FrameworkDescription = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.FrameworkTags = Tags(tags.IgnoreAWS())
	}

	log.Printf("[DEBUG] Creating Backup Framework: %#v", input)
	resp, err := conn.CreateFramework(input)
	if err != nil {
		return fmt.Errorf("error creating Backup Framework: %w", err)
	}

	// Set ID with the name since the name is unique for the framework
	d.SetId(aws.StringValue(resp.FrameworkName))

	// waiter since the status changes from CREATE_IN_PROGRESS to either COMPLETED or FAILED
	if _, err := waitFrameworkCreated(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("error waiting for Framework (%s) creation: %w", d.Id(), err)
	}

	return resourceFrameworkRead(d, meta)
}

func resourceFrameworkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	resp, err := conn.DescribeFramework(&backup.DescribeFrameworkInput{
		FrameworkName: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, backup.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Backup Framework (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading Backup Framework (%s): %w", d.Id(), err)
	}

	d.Set("arn", resp.FrameworkArn)
	d.Set("deployment_status", resp.DeploymentStatus)
	d.Set("description", resp.FrameworkDescription)
	d.Set("name", resp.FrameworkName)
	d.Set("status", resp.FrameworkStatus)

	if err := d.Set("creation_time", resp.CreationTime.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting creation_time: %s", err)
	}

	if err := d.Set("control", flattenFrameworkControls(resp.FrameworkControls)); err != nil {
		return fmt.Errorf("error setting control: %w", err)
	}

	tags, err := ListTags(conn, d.Get("arn").(string))
	if err != nil {
		return fmt.Errorf("error listing tags for Backup Framework (%s): %w", d.Id(), err)
	}
	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceFrameworkUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn()

	if d.HasChanges("description", "control") {
		input := &backup.UpdateFrameworkInput{
			IdempotencyToken:     aws.String(resource.UniqueId()),
			FrameworkControls:    expandFrameworkControls(d.Get("control").(*schema.Set).List()),
			FrameworkDescription: aws.String(d.Get("description").(string)),
			FrameworkName:        aws.String(d.Id()),
		}

		log.Printf("[DEBUG] Updating Backup Framework: %#v", input)

		_, err := tfresource.RetryWhenAWSErrCodeEquals(d.Timeout(schema.TimeoutUpdate), func() (interface{}, error) {
			return conn.UpdateFramework(input)
		}, backup.ErrCodeConflictException)

		if err != nil {
			return fmt.Errorf("error updating Backup Framework (%s): %w", d.Id(), err)
		}

		if _, err := waitFrameworkUpdated(conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for Framework (%s) update: %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags for Backup Framework (%s): %w", d.Id(), err)
		}
	}

	return resourceFrameworkRead(d, meta)
}

func resourceFrameworkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn()

	input := &backup.DeleteFrameworkInput{
		FrameworkName: aws.String(d.Id()),
	}

	_, err := tfresource.RetryWhenAWSErrCodeEquals(d.Timeout(schema.TimeoutDelete), func() (interface{}, error) {
		return conn.DeleteFramework(input)
	}, backup.ErrCodeConflictException)

	if err != nil {
		return fmt.Errorf("error deleting Backup Framework (%s): %w", d.Id(), err)
	}

	if _, err := waitFrameworkDeleted(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for Framework (%s) deletion: %w", d.Id(), err)
	}

	return nil
}

func expandFrameworkControls(controls []interface{}) []*backup.FrameworkControl {
	if len(controls) == 0 {
		return nil
	}

	frameworkControls := []*backup.FrameworkControl{}

	for _, control := range controls {
		tfMap := control.(map[string]interface{})

		// on some updates, there is an { ControlName: "" } element in Framework Controls.
		// this element must be skipped to avoid the "A control name is required." error
		// this happens for Step 7/7 for TestAccBackupFramework_updateControlScope
		if v, ok := tfMap["name"].(string); ok && v == "" {
			continue
		}

		frameworkControl := &backup.FrameworkControl{
			ControlName:  aws.String(tfMap["name"].(string)),
			ControlScope: expandControlScope(tfMap["scope"].([]interface{})),
		}

		if v, ok := tfMap["input_parameter"]; ok && v.(*schema.Set).Len() > 0 {
			frameworkControl.ControlInputParameters = expandInputParmaeters(tfMap["input_parameter"].(*schema.Set).List())
		}

		frameworkControls = append(frameworkControls, frameworkControl)
	}

	return frameworkControls
}

func expandInputParmaeters(inputParams []interface{}) []*backup.ControlInputParameter {
	if len(inputParams) == 0 {
		return nil
	}

	controlInputParameters := []*backup.ControlInputParameter{}

	for _, inputParam := range inputParams {
		tfMap := inputParam.(map[string]interface{})
		controlInputParameter := &backup.ControlInputParameter{}

		if v, ok := tfMap["name"].(string); ok && v != "" {
			controlInputParameter.ParameterName = aws.String(v)
		}

		if v, ok := tfMap["value"].(string); ok && v != "" {
			controlInputParameter.ParameterValue = aws.String(v)
		}

		controlInputParameters = append(controlInputParameters, controlInputParameter)
	}

	return controlInputParameters
}

func expandControlScope(scope []interface{}) *backup.ControlScope {
	if len(scope) == 0 || scope[0] == nil {
		return nil
	}

	tfMap, ok := scope[0].(map[string]interface{})
	if !ok {
		return nil
	}

	controlScope := &backup.ControlScope{}

	if v, ok := tfMap["compliance_resource_ids"]; ok && v.(*schema.Set).Len() > 0 {
		controlScope.ComplianceResourceIds = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := tfMap["compliance_resource_types"]; ok && v.(*schema.Set).Len() > 0 {
		controlScope.ComplianceResourceTypes = flex.ExpandStringSet(v.(*schema.Set))
	}

	// A maximum of one key-value pair can be provided.
	// The tag value is optional, but it cannot be an empty string
	if v, ok := tfMap["tags"].(map[string]interface{}); ok && len(v) > 0 {
		controlScope.Tags = Tags(tftags.New(v).IgnoreAWS())
	}

	return controlScope
}

func flattenFrameworkControls(controls []*backup.FrameworkControl) []interface{} {
	if controls == nil {
		return []interface{}{}
	}

	frameworkControls := []interface{}{}
	for _, control := range controls {
		values := map[string]interface{}{}
		values["input_parameter"] = flattenInputParameters(control.ControlInputParameters)
		values["name"] = aws.StringValue(control.ControlName)
		values["scope"] = flattenScope(control.ControlScope)
		frameworkControls = append(frameworkControls, values)
	}
	return frameworkControls
}

func flattenInputParameters(inputParams []*backup.ControlInputParameter) []interface{} {
	if inputParams == nil {
		return []interface{}{}
	}

	controlInputParameters := []interface{}{}
	for _, inputParam := range inputParams {
		values := map[string]interface{}{}
		values["name"] = aws.StringValue(inputParam.ParameterName)
		values["value"] = aws.StringValue(inputParam.ParameterValue)
		controlInputParameters = append(controlInputParameters, values)
	}
	return controlInputParameters
}

func flattenScope(scope *backup.ControlScope) []interface{} {
	if scope == nil {
		return []interface{}{}
	}

	controlScope := map[string]interface{}{
		"compliance_resource_ids":   flex.FlattenStringList(scope.ComplianceResourceIds),
		"compliance_resource_types": flex.FlattenStringList(scope.ComplianceResourceTypes),
	}

	if v := scope.Tags; v != nil {
		controlScope["tags"] = KeyValueTags(v).IgnoreAWS().Map()
	}

	return []interface{}{controlScope}
}
