package servicecatalog

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceProduct() *schema.Resource {
	return &schema.Resource{
		Create: resourceProductCreate,
		Read:   resourceProductRead,
		Update: resourceProductUpdate,
		Delete: resourceProductDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ProductReadyTimeout),
			Read:   schema.DefaultTimeout(ProductReadTimeout),
			Update: schema.DefaultTimeout(ProductUpdateTimeout),
			Delete: schema.DefaultTimeout(ProductDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      AcceptLanguageEnglish,
				ValidateFunc: validation.StringInSlice(AcceptLanguage_Values(), false),
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"distributor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"has_default_path": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provisioning_artifact_parameters": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disable_template_validation": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"template_physical_id": {
							Type:     schema.TypeString,
							Optional: true,
							ExactlyOneOf: []string{
								"provisioning_artifact_parameters.0.template_url",
								"provisioning_artifact_parameters.0.template_physical_id",
							},
						},
						"template_url": {
							Type:     schema.TypeString,
							Optional: true,
							ExactlyOneOf: []string{
								"provisioning_artifact_parameters.0.template_url",
								"provisioning_artifact_parameters.0.template_physical_id",
							},
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(servicecatalog.ProvisioningArtifactType_Values(), false),
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(servicecatalog.ProductType_Values(), false),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceProductCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &servicecatalog.CreateProductInput{
		IdempotencyToken: aws.String(resource.UniqueId()),
		Name:             aws.String(d.Get("name").(string)),
		Owner:            aws.String(d.Get("owner").(string)),
		ProductType:      aws.String(d.Get("type").(string)),
		ProvisioningArtifactParameters: expandProvisioningArtifactParameters(
			d.Get("provisioning_artifact_parameters").([]interface{})[0].(map[string]interface{}),
		),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("distributor"); ok {
		input.Distributor = aws.String(v.(string))
	}

	if v, ok := d.GetOk("support_description"); ok {
		input.SupportDescription = aws.String(v.(string))
	}

	if v, ok := d.GetOk("support_email"); ok {
		input.SupportEmail = aws.String(v.(string))
	}

	if v, ok := d.GetOk("support_url"); ok {
		input.SupportUrl = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	var output *servicecatalog.CreateProductOutput
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		output, err = conn.CreateProduct(input)

		if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateProduct(input)
	}

	if err != nil {
		return fmt.Errorf("error creating Service Catalog Product: %w", err)
	}

	if output == nil {
		return fmt.Errorf("error creating Service Catalog Product: empty response")
	}

	if output.ProductViewDetail == nil || output.ProductViewDetail.ProductViewSummary == nil {
		return fmt.Errorf("error creating Service Catalog Product: no product view detail or summary")
	}

	if output.ProvisioningArtifactDetail == nil {
		return fmt.Errorf("error creating Service Catalog Product: no provisioning artifact detail")
	}

	d.SetId(aws.StringValue(output.ProductViewDetail.ProductViewSummary.ProductId))

	if _, err := WaitProductReady(conn, aws.StringValue(input.AcceptLanguage),
		aws.StringValue(output.ProductViewDetail.ProductViewSummary.ProductId), d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("error waiting for Service Catalog Product (%s) to be ready: %w", d.Id(), err)
	}

	return resourceProductRead(d, meta)
}

func resourceProductRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	output, err := WaitProductReady(conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutRead))

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Service Catalog Product (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing Service Catalog Product (%s): %w", d.Id(), err)
	}

	if output == nil || output.ProductViewDetail == nil || output.ProductViewDetail.ProductViewSummary == nil {
		return fmt.Errorf("error getting Service Catalog Product (%s): empty response", d.Id())
	}

	pvs := output.ProductViewDetail.ProductViewSummary

	d.Set("arn", output.ProductViewDetail.ProductARN)
	if output.ProductViewDetail.CreatedTime != nil {
		d.Set("created_time", output.ProductViewDetail.CreatedTime.Format(time.RFC3339))
	}
	d.Set("description", pvs.ShortDescription)
	d.Set("distributor", pvs.Distributor)
	d.Set("has_default_path", pvs.HasDefaultPath)
	d.Set("name", pvs.Name)
	d.Set("owner", pvs.Owner)
	d.Set("status", output.ProductViewDetail.Status)
	d.Set("support_description", pvs.SupportDescription)
	d.Set("support_email", pvs.SupportEmail)
	d.Set("support_url", pvs.SupportUrl)
	d.Set("type", pvs.Type)

	tags := KeyValueTags(output.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceProductUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn()

	if d.HasChangesExcept("tags", "tags_all") {
		input := &servicecatalog.UpdateProductInput{
			Id: aws.String(d.Id()),
		}

		if v, ok := d.GetOk("accept_language"); ok {
			input.AcceptLanguage = aws.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			input.Description = aws.String(v.(string))
		}

		if v, ok := d.GetOk("distributor"); ok {
			input.Distributor = aws.String(v.(string))
		}

		if v, ok := d.GetOk("name"); ok {
			input.Name = aws.String(v.(string))
		}

		if v, ok := d.GetOk("owner"); ok {
			input.Owner = aws.String(v.(string))
		}

		if v, ok := d.GetOk("support_description"); ok {
			input.SupportDescription = aws.String(v.(string))
		}

		if v, ok := d.GetOk("support_email"); ok {
			input.SupportEmail = aws.String(v.(string))
		}

		if v, ok := d.GetOk("support_url"); ok {
			input.SupportUrl = aws.String(v.(string))
		}

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := conn.UpdateProduct(input)

			if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
				return resource.RetryableError(err)
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if tfresource.TimedOut(err) {
			_, err = conn.UpdateProduct(input)
		}

		if err != nil {
			return fmt.Errorf("error updating Service Catalog Product (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := productUpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating tags for Service Catalog Product (%s): %w", d.Id(), err)
		}
	}

	return resourceProductRead(d, meta)
}

func resourceProductDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn()

	input := &servicecatalog.DeleteProductInput{
		Id: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	_, err := conn.DeleteProduct(input)

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Service Catalog Product (%s): %w", d.Id(), err)
	}

	if _, err := WaitProductDeleted(conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for Service Catalog Product (%s) to be deleted: %w", d.Id(), err)
	}

	return nil
}

func expandProvisioningArtifactParameters(tfMap map[string]interface{}) *servicecatalog.ProvisioningArtifactProperties {
	if tfMap == nil {
		return nil
	}

	apiObject := &servicecatalog.ProvisioningArtifactProperties{}

	if v, ok := tfMap["description"].(string); ok && v != "" {
		apiObject.Description = aws.String(v)
	}

	if v, ok := tfMap["disable_template_validation"].(bool); ok {
		apiObject.DisableTemplateValidation = aws.Bool(v)
	}

	info := make(map[string]*string)

	// schema will enforce that one of these is present
	if v, ok := tfMap["template_physical_id"].(string); ok && v != "" {
		info["ImportFromPhysicalId"] = aws.String(v)
	}

	if v, ok := tfMap["template_url"].(string); ok && v != "" {
		info["LoadTemplateFromURL"] = aws.String(v)
	}

	apiObject.Info = info

	if v, ok := tfMap["name"].(string); ok && v != "" {
		apiObject.Name = aws.String(v)
	}

	if v, ok := tfMap["type"].(string); ok && v != "" {
		apiObject.Type = aws.String(v)
	}

	return apiObject
}
