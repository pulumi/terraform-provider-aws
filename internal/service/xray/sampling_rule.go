package xray

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceSamplingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSamplingRuleCreate,
		Read:   resourceSamplingRuleRead,
		Update: resourceSamplingRuleUpdate,
		Delete: resourceSamplingRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 32),
			},
			"resource_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 9999),
			},
			"fixed_rate": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"reservoir_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"service_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"host": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"http_method": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 10),
			},
			"url_path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"version": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(1, 32),
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},
	}
}

func resourceSamplingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	samplingRule := &xray.SamplingRule{
		RuleName:      aws.String(d.Get("rule_name").(string)),
		ResourceARN:   aws.String(d.Get("resource_arn").(string)),
		Priority:      aws.Int64(int64(d.Get("priority").(int))),
		FixedRate:     aws.Float64(d.Get("fixed_rate").(float64)),
		ReservoirSize: aws.Int64(int64(d.Get("reservoir_size").(int))),
		ServiceName:   aws.String(d.Get("service_name").(string)),
		ServiceType:   aws.String(d.Get("service_type").(string)),
		Host:          aws.String(d.Get("host").(string)),
		HTTPMethod:    aws.String(d.Get("http_method").(string)),
		URLPath:       aws.String(d.Get("url_path").(string)),
		Version:       aws.Int64(int64(d.Get("version").(int))),
	}

	if v, ok := d.GetOk("attributes"); ok {
		samplingRule.Attributes = flex.ExpandStringMap(v.(map[string]interface{}))
	}

	params := &xray.CreateSamplingRuleInput{
		SamplingRule: samplingRule,
		Tags:         Tags(tags.IgnoreAWS()),
	}

	out, err := conn.CreateSamplingRule(params)
	if err != nil {
		return fmt.Errorf("error creating XRay Sampling Rule: %w", err)
	}

	d.SetId(aws.StringValue(out.SamplingRuleRecord.SamplingRule.RuleName))

	return resourceSamplingRuleRead(d, meta)
}

func resourceSamplingRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	samplingRule, err := GetSamplingRule(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error reading XRay Sampling Rule (%s): %w", d.Id(), err)
	}

	if samplingRule == nil {
		log.Printf("[WARN] XRay Sampling Rule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	arn := aws.StringValue(samplingRule.RuleARN)
	d.Set("rule_name", samplingRule.RuleName)
	d.Set("resource_arn", samplingRule.ResourceARN)
	d.Set("priority", samplingRule.Priority)
	d.Set("fixed_rate", samplingRule.FixedRate)
	d.Set("reservoir_size", samplingRule.ReservoirSize)
	d.Set("service_name", samplingRule.ServiceName)
	d.Set("service_type", samplingRule.ServiceType)
	d.Set("host", samplingRule.Host)
	d.Set("http_method", samplingRule.HTTPMethod)
	d.Set("url_path", samplingRule.URLPath)
	d.Set("version", samplingRule.Version)
	d.Set("attributes", aws.StringValueMap(samplingRule.Attributes))
	d.Set("arn", arn)

	tags, err := ListTags(conn, arn)
	if err != nil {
		return fmt.Errorf("error listing tags for Xray Sampling group (%q): %s", d.Id(), err)
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

func resourceSamplingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %w", err)
		}
	}

	if d.HasChanges("attributes", "priority", "fixed_rate", "reservoir_size", "service_name", "service_type",
		"host", "http_method", "url_path", "resource_arn") {
		samplingRuleUpdate := &xray.SamplingRuleUpdate{
			RuleName:      aws.String(d.Id()),
			Priority:      aws.Int64(int64(d.Get("priority").(int))),
			FixedRate:     aws.Float64(d.Get("fixed_rate").(float64)),
			ReservoirSize: aws.Int64(int64(d.Get("reservoir_size").(int))),
			ServiceName:   aws.String(d.Get("service_name").(string)),
			ServiceType:   aws.String(d.Get("service_type").(string)),
			Host:          aws.String(d.Get("host").(string)),
			HTTPMethod:    aws.String(d.Get("http_method").(string)),
			URLPath:       aws.String(d.Get("url_path").(string)),
			ResourceARN:   aws.String(d.Get("resource_arn").(string)),
		}

		if d.HasChange("attributes") {
			attributes := map[string]*string{}
			if v, ok := d.GetOk("attributes"); ok {
				if m, ok := v.(map[string]interface{}); ok {
					attributes = flex.ExpandStringMap(m)
				}
			}
			samplingRuleUpdate.Attributes = attributes
		}

		params := &xray.UpdateSamplingRuleInput{
			SamplingRuleUpdate: samplingRuleUpdate,
		}

		_, err := conn.UpdateSamplingRule(params)
		if err != nil {
			return fmt.Errorf("error updating XRay Sampling Rule (%s): %w", d.Id(), err)
		}
	}

	return resourceSamplingRuleRead(d, meta)
}

func resourceSamplingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn

	log.Printf("[INFO] Deleting XRay Sampling Rule: %s", d.Id())

	params := &xray.DeleteSamplingRuleInput{
		RuleName: aws.String(d.Id()),
	}
	_, err := conn.DeleteSamplingRule(params)
	if err != nil {
		return fmt.Errorf("error deleting XRay Sampling Rule: %s", d.Id())
	}

	return nil
}

func GetSamplingRule(conn *xray.XRay, ruleName string) (*xray.SamplingRule, error) {
	params := &xray.GetSamplingRulesInput{}
	for {
		out, err := conn.GetSamplingRules(params)
		if err != nil {
			return nil, err
		}
		for _, samplingRuleRecord := range out.SamplingRuleRecords {
			samplingRule := samplingRuleRecord.SamplingRule
			if aws.StringValue(samplingRule.RuleName) == ruleName {
				return samplingRule, nil
			}
		}
		if aws.StringValue(out.NextToken) == "" {
			break
		}
		params.NextToken = out.NextToken
	}
	return nil, nil
}
