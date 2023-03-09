package docdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const clusterParameterGroupMaxParamsBulkEdit = 20

// @SDKResource("aws_docdb_cluster_parameter_group")
func ResourceClusterParameterGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceClusterParameterGroupCreate,
		ReadWithoutTimeout:   resourceClusterParameterGroupRead,
		UpdateWithoutTimeout: resourceClusterParameterGroupUpdate,
		DeleteWithoutTimeout: resourceClusterParameterGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Default:  "Managed by Pulumi",
			},
			"parameter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"apply_method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  docdb.ApplyMethodPendingReboot,
							ValidateFunc: validation.StringInSlice([]string{
								docdb.ApplyMethodImmediate,
								docdb.ApplyMethodPendingReboot,
							}, false),
						},
					},
				},
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceClusterParameterGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DocDBConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(ctx, d.Get("tags").(map[string]interface{})))

	var groupName string
	if v, ok := d.GetOk("name"); ok {
		groupName = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		groupName = resource.PrefixedUniqueId(v.(string))
	} else {
		groupName = resource.UniqueId()
	}

	createOpts := docdb.CreateDBClusterParameterGroupInput{
		DBClusterParameterGroupName: aws.String(groupName),
		DBParameterGroupFamily:      aws.String(d.Get("family").(string)),
		Description:                 aws.String(d.Get("description").(string)),
		Tags:                        Tags(tags.IgnoreAWS()),
	}

	log.Printf("[DEBUG] Create DocDB Cluster Parameter Group: %#v", createOpts)

	resp, err := conn.CreateDBClusterParameterGroupWithContext(ctx, &createOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "Error creating DocDB Cluster Parameter Group: %s", err)
	}

	d.SetId(aws.StringValue(createOpts.DBClusterParameterGroupName))

	d.Set("arn", resp.DBClusterParameterGroup.DBClusterParameterGroupArn)

	return append(diags, resourceClusterParameterGroupUpdate(ctx, d, meta)...)
}

func resourceClusterParameterGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DocDBConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	describeOpts := &docdb.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: aws.String(d.Id()),
	}

	describeResp, err := conn.DescribeDBClusterParameterGroupsWithContext(ctx, describeOpts)
	if err != nil {
		if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, docdb.ErrCodeDBParameterGroupNotFoundFault) {
			log.Printf("[WARN] DocDB Cluster Parameter Group (%s) not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "reading DocDB Cluster Parameter Group (%s): %s", d.Id(), err)
	}

	if len(describeResp.DBClusterParameterGroups) != 1 ||
		aws.StringValue(describeResp.DBClusterParameterGroups[0].DBClusterParameterGroupName) != d.Id() {
		return sdkdiag.AppendErrorf(diags, "Unable to find Cluster Parameter Group: %#v", describeResp.DBClusterParameterGroups)
	}

	arn := aws.StringValue(describeResp.DBClusterParameterGroups[0].DBClusterParameterGroupArn)
	d.Set("arn", arn)
	d.Set("description", describeResp.DBClusterParameterGroups[0].Description)
	d.Set("family", describeResp.DBClusterParameterGroups[0].DBParameterGroupFamily)
	d.Set("name", describeResp.DBClusterParameterGroups[0].DBClusterParameterGroupName)

	describeParametersOpts := &docdb.DescribeDBClusterParametersInput{
		DBClusterParameterGroupName: aws.String(d.Id()),
	}

	describeParametersResp, err := conn.DescribeDBClusterParametersWithContext(ctx, describeParametersOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DocDB Cluster Parameter Group (%s) parameters: %s", d.Id(), err)
	}

	if err := d.Set("parameter", flattenParameters(describeParametersResp.Parameters, d.Get("parameter").(*schema.Set).List())); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting docdb cluster parameter: %s", err)
	}

	tags, err := ListTags(ctx, conn, d.Get("arn").(string))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing tags for DocumentDB Cluster Parameter Group (%s): %s", d.Get("arn").(string), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags_all: %s", err)
	}

	return diags
}

func resourceClusterParameterGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DocDBConn()

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

		parameters := expandParameters(ns.Difference(os).List())
		if len(parameters) > 0 {
			// We can only modify 20 parameters at a time, so walk them until
			// we've got them all.
			for parameters != nil {
				var paramsToModify []*docdb.Parameter
				if len(parameters) <= clusterParameterGroupMaxParamsBulkEdit {
					paramsToModify, parameters = parameters[:], nil
				} else {
					paramsToModify, parameters = parameters[:clusterParameterGroupMaxParamsBulkEdit], parameters[clusterParameterGroupMaxParamsBulkEdit:]
				}
				parameterGroupName := d.Id()
				modifyOpts := docdb.ModifyDBClusterParameterGroupInput{
					DBClusterParameterGroupName: aws.String(parameterGroupName),
					Parameters:                  paramsToModify,
				}

				_, err := conn.ModifyDBClusterParameterGroupWithContext(ctx, &modifyOpts)
				if err != nil {
					return sdkdiag.AppendErrorf(diags, "Error modifying DocDB Cluster Parameter Group: %s", err)
				}
			}
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(ctx, conn, d.Get("arn").(string), o, n); err != nil {
			return sdkdiag.AppendErrorf(diags, "updating DocumentDB Cluster Parameter Group (%s) tags: %s", d.Get("arn").(string), err)
		}
	}

	return append(diags, resourceClusterParameterGroupRead(ctx, d, meta)...)
}

func resourceClusterParameterGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DocDBConn()

	deleteOpts := &docdb.DeleteDBClusterParameterGroupInput{
		DBClusterParameterGroupName: aws.String(d.Id()),
	}

	_, err := conn.DeleteDBClusterParameterGroupWithContext(ctx, deleteOpts)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, docdb.ErrCodeDBParameterGroupNotFoundFault) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "deleting DocumentDB Cluster Parameter Group (%s): %s", d.Id(), err)
	}

	if err := WaitForClusterParameterGroupDeletion(ctx, conn, d.Id()); err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting DocumentDB Cluster Parameter Group (%s): %s", d.Id(), err)
	}
	return diags
}

func WaitForClusterParameterGroupDeletion(ctx context.Context, conn *docdb.DocDB, name string) error {
	params := &docdb.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: aws.String(name),
	}

	err := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		_, err := conn.DescribeDBClusterParameterGroupsWithContext(ctx, params)

		if tfawserr.ErrCodeEquals(err, docdb.ErrCodeDBParameterGroupNotFoundFault) {
			return nil
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("DocDB Parameter Group (%s) still exists", name))
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DescribeDBClusterParameterGroupsWithContext(ctx, params)
		if tfawserr.ErrCodeEquals(err, docdb.ErrCodeDBParameterGroupNotFoundFault) {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf("Error deleting DocDB cluster parameter group: %s", err)
	}
	return nil
}
