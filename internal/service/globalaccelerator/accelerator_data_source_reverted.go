package globalaccelerator

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceAccelerator() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAcceleratorRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_address_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosted_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ip_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flow_logs_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"flow_logs_s3_bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_logs_s3_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

// Global Route53 Zone ID for Global Accelerators, exported as a
// convenience attribute for Route53 aliases (see
// https://docs.aws.amazon.com/Route53/latest/APIReference/API_AliasTarget.html).
const route53ZoneID = "Z2BJ6XQ5FK7U4H"

func dataSourceAcceleratorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).GlobalAcceleratorConn()
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	var results []*globalaccelerator.Accelerator

	err := conn.ListAcceleratorsPages(&globalaccelerator.ListAcceleratorsInput{}, func(page *globalaccelerator.ListAcceleratorsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, l := range page.Accelerators {
			if l == nil {
				continue
			}

			if v, ok := d.GetOk("arn"); ok && v.(string) != aws.StringValue(l.AcceleratorArn) {
				continue
			}

			if v, ok := d.GetOk("name"); ok && v.(string) != aws.StringValue(l.Name) {
				continue
			}

			results = append(results, l)
		}

		return !lastPage
	})

	if err != nil {
		return diag.Errorf("error reading AWS Global Accelerator: %w", err)
	}

	if len(results) != 1 {
		return diag.Errorf("Search returned %d results, please revise so only one is returned", len(results))
	}

	accelerator := results[0]
	d.SetId(aws.StringValue(accelerator.AcceleratorArn))
	d.Set("arn", accelerator.AcceleratorArn)
	d.Set("enabled", accelerator.Enabled)
	d.Set("dns_name", accelerator.DnsName)
	d.Set("hosted_zone_id", route53ZoneID)
	d.Set("name", accelerator.Name)
	d.Set("ip_address_type", accelerator.IpAddressType)
	d.Set("ip_sets", flattenIPSets(accelerator.IpSets))

	acceleratorAttributes, err := FindAcceleratorAttributesByARN(ctx, conn, d.Id())
	if err != nil {
		return diag.Errorf("error reading Global Accelerator Accelerator (%s) attributes: %w", d.Id(), err)
	}

	if err := d.Set("attributes", []interface{}{flattenAcceleratorAttributes(acceleratorAttributes)}); err != nil {
		return diag.Errorf("error setting attributes: %w", err)
	}

	tags, err := ListTags(ctx, conn, d.Id())
	if err != nil {
		return diag.Errorf("error listing tags for Global Accelerator Accelerator (%s): %w", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return diag.Errorf("error setting tags: %w", err)
	}
	return nil
}
