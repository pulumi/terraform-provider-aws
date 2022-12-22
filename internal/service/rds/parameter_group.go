package rds

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	rds_sdkv2 "github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create:               resourceParameterGroupCreate,
		Read:                 resourceParameterGroupRead,
		Update:               resourceParameterGroupUpdate,
		DeleteWithoutTimeout: resourceParameterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validParamGroupName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validParamGroupNamePrefix,
			},
			"family": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Managed by Terraform",
			},
			"parameter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apply_method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "immediate",
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceParameterHash,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	var groupName string
	if v, ok := d.GetOk("name"); ok {
		groupName = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		groupName = resource.PrefixedUniqueId(v.(string))
	} else {
		groupName = resource.UniqueId()
	}
	d.Set("name", groupName)

	createOpts := rds.CreateDBParameterGroupInput{
		DBParameterGroupName:   aws.String(groupName),
		DBParameterGroupFamily: aws.String(d.Get("family").(string)),
		Description:            aws.String(d.Get("description").(string)),
		Tags:                   Tags(tags.IgnoreAWS()),
	}

	log.Printf("[DEBUG] Create DB Parameter Group: %#v", createOpts)
	resp, err := conn.CreateDBParameterGroup(&createOpts)
	if err != nil {
		return fmt.Errorf("creating DB Parameter Group: %s", err)
	}

	d.SetId(aws.StringValue(resp.DBParameterGroup.DBParameterGroupName))
	d.Set("arn", resp.DBParameterGroup.DBParameterGroupArn)
	log.Printf("[INFO] DB Parameter Group ID: %s", d.Id())

	return resourceParameterGroupUpdate(d, meta)
}

func resourceParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	describeOpts := rds.DescribeDBParameterGroupsInput{
		DBParameterGroupName: aws.String(d.Id()),
	}

	describeResp, err := conn.DescribeDBParameterGroups(&describeOpts)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBParameterGroupNotFoundFault) {
			log.Printf("[WARN] DB Parameter Group (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if len(describeResp.DBParameterGroups) != 1 ||
		aws.StringValue(describeResp.DBParameterGroups[0].DBParameterGroupName) != d.Id() {
		return fmt.Errorf("Unable to find Parameter Group: %#v", describeResp.DBParameterGroups)
	}

	d.Set("name", describeResp.DBParameterGroups[0].DBParameterGroupName)
	d.Set("family", describeResp.DBParameterGroups[0].DBParameterGroupFamily)
	d.Set("description", describeResp.DBParameterGroups[0].Description)

	configParams := d.Get("parameter").(*schema.Set)
	describeParametersOpts := rds.DescribeDBParametersInput{
		DBParameterGroupName: aws.String(d.Id()),
	}
	if configParams.Len() < 1 {
		// if we don't have any params in the ResourceData already, two possibilities
		// first, we don't have a config available to us. Second, we do, but it has
		// no parameters. We're going to assume the first, to be safe. In this case,
		// we're only going to ask for the user-modified values, because any defaults
		// the user may have _also_ set are indistinguishable from the hundreds of
		// defaults AWS sets. If the user hasn't set any parameters, this will return
		// an empty list anyways, so we just make some unnecessary requests. But in
		// the more common case (I assume) of an import, this will make fewer requests
		// and "do the right thing".
		describeParametersOpts.Source = aws.String("user")
	}

	var parameters []*rds.Parameter
	err = conn.DescribeDBParametersPages(&describeParametersOpts,
		func(describeParametersResp *rds.DescribeDBParametersOutput, lastPage bool) bool {
			parameters = append(parameters, describeParametersResp.Parameters...)
			return !lastPage
		})
	if err != nil {
		return err
	}

	var userParams []*rds.Parameter
	if configParams.Len() < 1 {
		// if we have no config/no parameters in config, we've already asked for only
		// user-modified values, so we can just use the entire response.
		userParams = parameters
	} else {
		// if we have a config available to us, we have two possible classes of value
		// in the config. On the one hand, the user could have specified a parameter
		// that _actually_ changed things, in which case its Source would be set to
		// user. On the other, they may have specified a parameter that coincides with
		// the default value. In that case, the Source will be set to "system" or
		// "engine-default". We need to set the union of all "user" Source parameters
		// _and_ the "system"/"engine-default" Source parameters _that appear in the
		// config_ in the state, or the user gets a perpetual diff. See
		// terraform-providers/terraform-provider-aws#593 for more context and details.
		confParams := expandParameters(configParams.List())
		for _, param := range parameters {
			if param.Source == nil || param.ParameterName == nil {
				continue
			}
			if aws.StringValue(param.Source) == "user" {
				userParams = append(userParams, param)
				continue
			}
			var paramFound bool
			for _, cp := range confParams {
				if cp.ParameterName == nil {
					continue
				}
				if aws.StringValue(cp.ParameterName) == aws.StringValue(param.ParameterName) {
					userParams = append(userParams, param)
					break
				}
			}
			if !paramFound {
				log.Printf("[DEBUG] Not persisting %s to state, as its source is %q and it isn't in the config", aws.StringValue(param.ParameterName), aws.StringValue(param.Source))
			}
		}
	}

	err = d.Set("parameter", flattenParameters(userParams))
	if err != nil {
		return fmt.Errorf("setting 'parameter' in state: %w", err)
	}

	arn := aws.StringValue(describeResp.DBParameterGroups[0].DBParameterGroupArn)
	d.Set("arn", arn)

	tags, err := ListTags(conn, d.Get("arn").(string))

	if err != nil {
		return fmt.Errorf("listing tags for RDS DB Parameter Group (%s): %w", d.Get("arn").(string), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	return nil
}

const maxParamModifyChunk = 20

func resourceParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	if d.HasChange("parameter") {
		o, n := d.GetChange("parameter")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		// Expand the "parameter" set to aws-sdk-go compat []rds.Parameter
		parameters := expandParameters(ns.Difference(os).List())

		if len(parameters) > 0 {
			// We can only modify 20 parameters at a time, so walk them until
			// we've got them all.

			for parameters != nil {
				var paramsToModify []*rds.Parameter
				paramsToModify, parameters = ResourceParameterModifyChunk(parameters, maxParamModifyChunk)

				modifyOpts := rds.ModifyDBParameterGroupInput{
					DBParameterGroupName: aws.String(d.Get("name").(string)),
					Parameters:           paramsToModify,
				}

				log.Printf("[DEBUG] Modify DB Parameter Group: %s", modifyOpts)
				_, err := conn.ModifyDBParameterGroup(&modifyOpts)
				if err != nil {
					return fmt.Errorf("modifying DB Parameter Group: %s", err)
				}
			}
		}

		toRemove := map[string]*rds.Parameter{}

		for _, p := range expandParameters(os.List()) {
			if p.ParameterName != nil {
				toRemove[*p.ParameterName] = p
			}
		}

		for _, p := range expandParameters(ns.List()) {
			if p.ParameterName != nil {
				delete(toRemove, *p.ParameterName)
			}
		}

		// Reset parameters that have been removed
		var resetParameters []*rds.Parameter
		for _, v := range toRemove {
			resetParameters = append(resetParameters, v)
		}
		if len(resetParameters) > 0 {
			for resetParameters != nil {
				var paramsToReset []*rds.Parameter
				if len(resetParameters) <= maxParamModifyChunk {
					paramsToReset, resetParameters = resetParameters[:], nil
				} else {
					paramsToReset, resetParameters = resetParameters[:maxParamModifyChunk], resetParameters[maxParamModifyChunk:]
				}

				parameterGroupName := d.Get("name").(string)
				resetOpts := rds.ResetDBParameterGroupInput{
					DBParameterGroupName: aws.String(parameterGroupName),
					Parameters:           paramsToReset,
					ResetAllParameters:   aws.Bool(false),
				}

				log.Printf("[DEBUG] Reset DB Parameter Group: %s", resetOpts)
				_, err := conn.ResetDBParameterGroup(&resetOpts)
				if err != nil {
					return fmt.Errorf("resetting DB Parameter Group: %s", err)
				}
			}
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("updating RDS DB Parameter Group (%s) tags: %s", d.Get("arn").(string), err)
		}
	}

	return resourceParameterGroupRead(d, meta)
}

func resourceParameterGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	conn := meta.(*conns.AWSClient).RDSClient()
	deleteOpts := rds_sdkv2.DeleteDBParameterGroupInput{
		DBParameterGroupName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting RDS DB Parameter Group: %s", d.Id())
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteDBParameterGroup(ctx, &deleteOpts)
		if errs.IsA[*types.DBParameterGroupNotFoundFault](err) {
			return nil
		} else if errs.IsA[*types.InvalidDBParameterGroupStateFault](err) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DeleteDBParameterGroup(ctx, &deleteOpts)
	}
	if err != nil {
		return errs.AppendErrorf(diags, "deleting RDS DB Parameter Group (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceParameterHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	// Store the value as a lower case string, to match how we store them in FlattenParameters
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["name"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["apply_method"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", m["value"].(string)))

	// This hash randomly affects the "order" of the set, which affects in what order parameters
	// are applied, when there are more than 20 (chunked).
	return create.StringHashcode(buf.String())
}

func ResourceParameterModifyChunk(all []*rds.Parameter, maxChunkSize int) ([]*rds.Parameter, []*rds.Parameter) {
	// Since the hash randomly affect the set "order," this attempts to prioritize important
	// parameters to go in the first chunk (i.e., charset)

	if len(all) <= maxChunkSize {
		return all[:], nil
	}

	var modifyChunk, remainder []*rds.Parameter

	// pass 1
	for i, p := range all {
		if len(modifyChunk) >= maxChunkSize {
			remainder = append(remainder, all[i:]...)
			return modifyChunk, remainder
		}

		if strings.Contains(aws.StringValue(p.ParameterName), "character_set") && aws.StringValue(p.ApplyMethod) != "pending-reboot" {
			modifyChunk = append(modifyChunk, p)
			continue
		}

		remainder = append(remainder, p)
	}

	all = remainder
	remainder = nil

	// pass 2 - avoid pending reboot
	for i, p := range all {
		if len(modifyChunk) >= maxChunkSize {
			remainder = append(remainder, all[i:]...)
			return modifyChunk, remainder
		}

		if aws.StringValue(p.ApplyMethod) != "pending-reboot" {
			modifyChunk = append(modifyChunk, p)
			continue
		}

		remainder = append(remainder, p)
	}

	all = remainder
	remainder = nil

	// pass 3 - everything else
	for i, p := range all {
		if len(modifyChunk) >= maxChunkSize {
			remainder = append(remainder, all[i:]...)
			return modifyChunk, remainder
		}

		modifyChunk = append(modifyChunk, p)
	}

	return modifyChunk, remainder
}
