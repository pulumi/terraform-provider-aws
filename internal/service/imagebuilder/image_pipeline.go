package imagebuilder

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceImagePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceImagePipelineCreate,
		Read:   resourceImagePipelineRead,
		Update: resourceImagePipelineUpdate,
		Delete: resourceImagePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_recipe_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^arn:aws[^:]*:imagebuilder:[^:]+:(?:\d{12}|aws):container-recipe/[a-z0-9-_]+/\d+\.\d+\.\d+$`), "valid container recipe ARN must be provided"),
				ExactlyOneOf: []string{"container_recipe_arn", "image_recipe_arn"},
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_last_run": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_next_run": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"distribution_configuration_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^arn:aws[^:]*:imagebuilder:[^:]+:(?:\d{12}|aws):distribution-configuration/[a-z0-9-_]+$`), "valid distribution configuration ARN must be provided"),
			},
			"enhanced_image_metadata_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"image_recipe_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^arn:aws[^:]*:imagebuilder:[^:]+:(?:\d{12}|aws):image-recipe/[a-z0-9-_]+/\d+\.\d+\.\d+$`), "valid image recipe ARN must be provided"),
				ExactlyOneOf: []string{"container_recipe_arn", "image_recipe_arn"},
			},
			"image_tests_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_tests_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"timeout_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      720,
							ValidateFunc: validation.IntBetween(60, 1440),
						},
					},
				},
			},
			"infrastructure_configuration_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^arn:aws[^:]*:imagebuilder:[^:]+:(?:\d{12}|aws):infrastructure-configuration/[a-z0-9-_]+$`), "valid infrastructure configuration ARN must be provided"),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[-_A-Za-z-0-9][-_A-Za-z0-9 ]{1,126}[-_A-Za-z-0-9]$"), "valid name must be provided"),
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_execution_start_condition": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      imagebuilder.PipelineExecutionStartConditionExpressionMatchAndDependencyUpdatesAvailable,
							ValidateFunc: validation.StringInSlice(imagebuilder.PipelineExecutionStartCondition_Values(), false),
						},
						"schedule_expression": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(3, 100),
								validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]{2,}(?:\/[a-zA-z0-9-_+]+)*`), "")),
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      imagebuilder.PipelineStatusEnabled,
				ValidateFunc: validation.StringInSlice(imagebuilder.PipelineStatus_Values(), false),
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceImagePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &imagebuilder.CreateImagePipelineInput{
		ClientToken:                  aws.String(resource.UniqueId()),
		EnhancedImageMetadataEnabled: aws.Bool(d.Get("enhanced_image_metadata_enabled").(bool)),
	}

	if v, ok := d.GetOk("container_recipe_arn"); ok {
		input.ContainerRecipeArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("distribution_configuration_arn"); ok {
		input.DistributionConfigurationArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("image_recipe_arn"); ok {
		input.ImageRecipeArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("image_tests_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.ImageTestsConfiguration = expandImageTestConfiguration(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("infrastructure_configuration_arn"); ok {
		input.InfrastructureConfigurationArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		input.Name = aws.String(v.(string))
	}

	if v, ok := d.GetOk("schedule"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.Schedule = expandPipelineSchedule(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("status"); ok {
		input.Status = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	output, err := conn.CreateImagePipeline(input)

	if err != nil {
		return fmt.Errorf("error creating Image Builder Image Pipeline: %w", err)
	}

	if output == nil {
		return fmt.Errorf("error creating Image Builder Image Pipeline: empty response")
	}

	d.SetId(aws.StringValue(output.ImagePipelineArn))

	return resourceImagePipelineRead(d, meta)
}

func resourceImagePipelineRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &imagebuilder.GetImagePipelineInput{
		ImagePipelineArn: aws.String(d.Id()),
	}

	output, err := conn.GetImagePipeline(input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, imagebuilder.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Image Builder Image Pipeline (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting Image Builder Image Pipeline (%s): %w", d.Id(), err)
	}

	if output == nil || output.ImagePipeline == nil {
		return fmt.Errorf("error getting Image Builder Image Pipeline (%s): empty response", d.Id())
	}

	imagePipeline := output.ImagePipeline

	d.Set("arn", imagePipeline.Arn)
	d.Set("container_recipe_arn", imagePipeline.ContainerRecipeArn)
	d.Set("date_created", imagePipeline.DateCreated)
	d.Set("date_last_run", imagePipeline.DateLastRun)
	d.Set("date_next_run", imagePipeline.DateNextRun)
	d.Set("date_updated", imagePipeline.DateUpdated)
	d.Set("description", imagePipeline.Description)
	d.Set("distribution_configuration_arn", imagePipeline.DistributionConfigurationArn)
	d.Set("enhanced_image_metadata_enabled", imagePipeline.EnhancedImageMetadataEnabled)
	d.Set("image_recipe_arn", imagePipeline.ImageRecipeArn)

	if imagePipeline.ImageTestsConfiguration != nil {
		d.Set("image_tests_configuration", []interface{}{flattenImageTestsConfiguration(imagePipeline.ImageTestsConfiguration)})
	} else {
		d.Set("image_tests_configuration", nil)
	}

	d.Set("infrastructure_configuration_arn", imagePipeline.InfrastructureConfigurationArn)
	d.Set("name", imagePipeline.Name)
	d.Set("platform", imagePipeline.Platform)

	if imagePipeline.Schedule != nil {
		d.Set("schedule", []interface{}{flattenSchedule(imagePipeline.Schedule)})
	} else {
		d.Set("schedule", nil)
	}

	d.Set("status", imagePipeline.Status)

	tags := KeyValueTags(imagePipeline.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceImagePipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn

	if d.HasChanges(
		"description",
		"distribution_configuration_arn",
		"enhanced_image_metadata_enabled",
		"image_tests_configuration",
		"infrastructure_configuration_arn",
		"schedule",
		"status",
	) {
		input := &imagebuilder.UpdateImagePipelineInput{
			ClientToken:                  aws.String(resource.UniqueId()),
			EnhancedImageMetadataEnabled: aws.Bool(d.Get("enhanced_image_metadata_enabled").(bool)),
			ImagePipelineArn:             aws.String(d.Id()),
		}

		if v, ok := d.GetOk("container_recipe_arn"); ok {
			input.ContainerRecipeArn = aws.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			input.Description = aws.String(v.(string))
		}

		if v, ok := d.GetOk("distribution_configuration_arn"); ok {
			input.DistributionConfigurationArn = aws.String(v.(string))
		}

		if v, ok := d.GetOk("image_recipe_arn"); ok {
			input.ImageRecipeArn = aws.String(v.(string))
		}

		if v, ok := d.GetOk("image_tests_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input.ImageTestsConfiguration = expandImageTestConfiguration(v.([]interface{})[0].(map[string]interface{}))
		}

		if v, ok := d.GetOk("infrastructure_configuration_arn"); ok {
			input.InfrastructureConfigurationArn = aws.String(v.(string))
		}

		if v, ok := d.GetOk("schedule"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input.Schedule = expandPipelineSchedule(v.([]interface{})[0].(map[string]interface{}))
		}

		if v, ok := d.GetOk("status"); ok {
			input.Status = aws.String(v.(string))
		}

		_, err := conn.UpdateImagePipeline(input)

		if err != nil {
			return fmt.Errorf("error updating Image Builder Image Pipeline (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating tags for Image Builder Image Pipeline (%s): %w", d.Id(), err)
		}
	}

	return resourceImagePipelineRead(d, meta)
}

func resourceImagePipelineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn

	input := &imagebuilder.DeleteImagePipelineInput{
		ImagePipelineArn: aws.String(d.Id()),
	}

	_, err := conn.DeleteImagePipeline(input)

	if tfawserr.ErrCodeEquals(err, imagebuilder.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Image Builder Image Pipeline (%s): %w", d.Id(), err)
	}

	return nil
}

func expandImageTestConfiguration(tfMap map[string]interface{}) *imagebuilder.ImageTestsConfiguration {
	if tfMap == nil {
		return nil
	}

	apiObject := &imagebuilder.ImageTestsConfiguration{}

	if v, ok := tfMap["image_tests_enabled"].(bool); ok {
		apiObject.ImageTestsEnabled = aws.Bool(v)
	}

	if v, ok := tfMap["timeout_minutes"].(int); ok && v != 0 {
		apiObject.TimeoutMinutes = aws.Int64(int64(v))
	}

	return apiObject
}

func expandPipelineSchedule(tfMap map[string]interface{}) *imagebuilder.Schedule {
	if tfMap == nil {
		return nil
	}

	apiObject := &imagebuilder.Schedule{}

	if v, ok := tfMap["pipeline_execution_start_condition"].(string); ok && v != "" {
		apiObject.PipelineExecutionStartCondition = aws.String(v)
	}

	if v, ok := tfMap["schedule_expression"].(string); ok && v != "" {
		apiObject.ScheduleExpression = aws.String(v)
	}

	if v, ok := tfMap["timezone"].(string); ok && v != "" {
		apiObject.Timezone = aws.String(v)
	}

	return apiObject
}

func flattenImageTestsConfiguration(apiObject *imagebuilder.ImageTestsConfiguration) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.ImageTestsEnabled; v != nil {
		tfMap["image_tests_enabled"] = aws.BoolValue(v)
	}

	if v := apiObject.TimeoutMinutes; v != nil {
		tfMap["timeout_minutes"] = aws.Int64Value(v)
	}

	return tfMap
}

func flattenSchedule(apiObject *imagebuilder.Schedule) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.PipelineExecutionStartCondition; v != nil {
		tfMap["pipeline_execution_start_condition"] = aws.StringValue(v)
	}

	if v := apiObject.ScheduleExpression; v != nil {
		tfMap["schedule_expression"] = aws.StringValue(v)
	}

	if v := apiObject.Timezone; v != nil {
		tfMap["timezone"] = aws.StringValue(v)
	}

	return tfMap
}
