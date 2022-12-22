package apprunner

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceAutoScalingConfigurationVersion() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAutoScalingConfigurationCreate,
		ReadWithoutTimeout:   resourceAutoScalingConfigurationRead,
		UpdateWithoutTimeout: resourceAutoScalingConfigurationUpdate,
		DeleteWithoutTimeout: resourceAutoScalingConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_scaling_configuration_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_scaling_configuration_revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"latest": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_concurrency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 200),
			},
			"max_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      25,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 25),
			},
			"min_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 25),
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

func resourceAutoScalingConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppRunnerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("auto_scaling_configuration_name").(string)

	input := &apprunner.CreateAutoScalingConfigurationInput{
		AutoScalingConfigurationName: aws.String(name),
	}

	if v, ok := d.GetOk("max_concurrency"); ok {
		input.MaxConcurrency = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("max_size"); ok {
		input.MaxSize = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("min_size"); ok {
		input.MinSize = aws.Int64(int64(v.(int)))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	output, err := conn.CreateAutoScalingConfigurationWithContext(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating App Runner AutoScaling Configuration Version (%s): %w", name, err))
	}

	if output == nil || output.AutoScalingConfiguration == nil {
		return diag.FromErr(fmt.Errorf("error creating App Runner AutoScaling Configuration Version (%s): empty output", name))
	}

	d.SetId(aws.StringValue(output.AutoScalingConfiguration.AutoScalingConfigurationArn))

	if err := WaitAutoScalingConfigurationActive(ctx, conn, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for AutoScaling Configuration Version (%s) creation: %w", d.Id(), err))
	}

	return resourceAutoScalingConfigurationRead(ctx, d, meta)
}

func resourceAutoScalingConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppRunnerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &apprunner.DescribeAutoScalingConfigurationInput{
		AutoScalingConfigurationArn: aws.String(d.Id()),
	}

	output, err := conn.DescribeAutoScalingConfigurationWithContext(ctx, input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, apprunner.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] App Runner AutoScaling Configuration Version (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading App Runner AutoScaling Configuration Version (%s): %w", d.Id(), err))
	}

	if output == nil || output.AutoScalingConfiguration == nil {
		return diag.FromErr(fmt.Errorf("error reading App Runner AutoScaling Configuration Version (%s): empty output", d.Id()))
	}

	if aws.StringValue(output.AutoScalingConfiguration.Status) == AutoScalingConfigurationStatusInactive {
		if d.IsNewResource() {
			return diag.FromErr(fmt.Errorf("error reading App Runner AutoScaling Configuration Version (%s): %s after creation", d.Id(), aws.StringValue(output.AutoScalingConfiguration.Status)))
		}
		log.Printf("[WARN] App Runner AutoScaling Configuration Version (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	config := output.AutoScalingConfiguration
	arn := aws.StringValue(config.AutoScalingConfigurationArn)

	d.Set("arn", arn)
	d.Set("auto_scaling_configuration_name", config.AutoScalingConfigurationName)
	d.Set("auto_scaling_configuration_revision", config.AutoScalingConfigurationRevision)
	d.Set("latest", config.Latest)
	d.Set("max_concurrency", config.MaxConcurrency)
	d.Set("max_size", config.MaxSize)
	d.Set("min_size", config.MinSize)
	d.Set("status", config.Status)

	tags, err := ListTags(conn, arn)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing tags for App Runner AutoScaling Configuration Version (%s): %s", arn, err))
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tags_all: %w", err))
	}

	return nil
}

func resourceAutoScalingConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppRunnerConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return diag.FromErr(fmt.Errorf("error updating App Runner AutoScaling Configuration Version (%s) tags: %s", d.Get("arn").(string), err))
		}
	}

	return resourceAutoScalingConfigurationRead(ctx, d, meta)
}

func resourceAutoScalingConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppRunnerConn

	input := &apprunner.DeleteAutoScalingConfigurationInput{
		AutoScalingConfigurationArn: aws.String(d.Id()),
	}

	_, err := conn.DeleteAutoScalingConfigurationWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, apprunner.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting App Runner AutoScaling Configuration Version (%s): %w", d.Id(), err))
	}

	if err := WaitAutoScalingConfigurationInactive(ctx, conn, d.Id()); err != nil {
		if tfawserr.ErrCodeEquals(err, apprunner.ErrCodeResourceNotFoundException) {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error waiting for AutoScaling Configuration Version (%s) deletion: %w", d.Id(), err))
	}

	return nil
}
