package route53recoveryreadiness

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53recoveryreadiness"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceCell() *schema.Resource {
	return &schema.Resource{
		Create: resourceCellCreate,
		Read:   resourceCellRead,
		Update: resourceCellUpdate,
		Delete: resourceCellDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cell_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cells": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parent_readiness_scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceCellCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53RecoveryReadinessConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &route53recoveryreadiness.CreateCellInput{
		CellName: aws.String(d.Get("cell_name").(string)),
		Cells:    flex.ExpandStringList(d.Get("cells").([]interface{})),
	}

	resp, err := conn.CreateCell(input)
	if err != nil {
		return fmt.Errorf("error creating Route53 Recovery Readiness Cell: %w", err)
	}

	d.SetId(aws.StringValue(resp.CellName))

	if len(tags) > 0 {
		arn := aws.StringValue(resp.CellArn)
		if err := UpdateTags(conn, arn, nil, tags); err != nil {
			return fmt.Errorf("error adding Route53 Recovery Readiness Cell (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceCellRead(d, meta)
}

func resourceCellRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53RecoveryReadinessConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &route53recoveryreadiness.GetCellInput{
		CellName: aws.String(d.Id()),
	}

	resp, err := conn.GetCell(input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, route53recoveryreadiness.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Route53RecoveryReadiness Cell (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing Route53 Recovery Readiness Cell: %s", err)
	}

	d.Set("arn", resp.CellArn)
	d.Set("cell_name", resp.CellName)
	d.Set("cells", resp.Cells)
	d.Set("parent_readiness_scopes", resp.ParentReadinessScopes)

	tags, err := ListTags(conn, d.Get("arn").(string))

	if err != nil {
		return fmt.Errorf("error listing tags for Route53 Recovery Readiness Cell (%s): %w", d.Id(), err)
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

func resourceCellUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53RecoveryReadinessConn

	input := &route53recoveryreadiness.UpdateCellInput{
		CellName: aws.String(d.Id()),
		Cells:    flex.ExpandStringList(d.Get("cells").([]interface{})),
	}

	_, err := conn.UpdateCell(input)
	if err != nil {
		return fmt.Errorf("error updating Route53 Recovery Readiness Cell: %s", err)
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		arn := d.Get("arn").(string)
		if err := UpdateTags(conn, arn, o, n); err != nil {
			return fmt.Errorf("error updating Route53 Recovery Readiness Cell (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceCellRead(d, meta)
}

func resourceCellDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53RecoveryReadinessConn

	input := &route53recoveryreadiness.DeleteCellInput{
		CellName: aws.String(d.Id()),
	}
	_, err := conn.DeleteCell(input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, route53recoveryreadiness.ErrCodeResourceNotFoundException) {
			return nil
		}
		return fmt.Errorf("error deleting Route53 Recovery Readiness Cell: %s", err)
	}

	gcinput := &route53recoveryreadiness.GetCellInput{
		CellName: aws.String(d.Id()),
	}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.GetCell(gcinput)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, route53recoveryreadiness.ErrCodeResourceNotFoundException) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("Route 53 Recovery Readiness Cell (%s) still exists", d.Id()))
	})
	if tfresource.TimedOut(err) {
		_, err = conn.GetCell(gcinput)
	}
	if err != nil {
		return fmt.Errorf("error waiting for Route 53 Recovery Readiness Cell (%s) deletion: %s", d.Id(), err)
	}

	return nil
}
