package aws

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func dataSourceAwsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsAwsDedicatedHostRead,

		Schema: map[string]*schema.Schema{
			"tags": tagsSchemaComputed(),
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_recovery": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_placement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsAwsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig
	conn := meta.(*AWSClient).ec2conn
	hostID, hostIDOk := d.GetOk("host_id")

	params := &ec2.DescribeHostsInput{}
	if hostIDOk {
		params.HostIds = []*string{aws.String(hostID.(string))}
	}
	resp, err := conn.DescribeHosts(params)
	if err != nil {
		return err
	}
	// If no hosts were returned, return
	if resp.Hosts == nil || len(resp.Hosts) == 0 {
		return errors.New("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(resp.Hosts) > 1 {
		return errors.New(`Your query returned more than one result. Please try a more 
			specific search criteria.`)
	}

	host := resp.Hosts[0]
	log.Printf("[DEBUG] aws_dedicated_host - Single host ID found: %s", *host.HostId)

	d.SetId(*host.HostId)
	d.Set("instance_state", host.State)

	if host.AvailabilityZone != nil {
		d.Set("availability_zone", host.AvailabilityZone)
	}
	if host.HostRecovery != nil {
		d.Set("host_recovery", host.HostRecovery)
	}
	if host.AutoPlacement != nil {
		d.Set("auto_placement", host.AutoPlacement)
	}
	if host.HostProperties.InstanceType != nil {
		d.Set("instance_type", host.HostProperties.InstanceType)
	}
	if host.HostProperties.InstanceFamily != nil {
		d.Set("instance_family", host.HostProperties.InstanceFamily)
	}
	if host.HostProperties.Cores != nil {
		d.Set("cores", host.HostProperties.Cores)
	}
	if host.HostProperties.Sockets != nil {
		d.Set("sockets", host.HostProperties.Sockets)
	}
	if host.HostProperties.TotalVCpus != nil {
		d.Set("total_vcpus", host.HostProperties.TotalVCpus)
	}

	if err := d.Set("tags", keyvaluetags.Ec2KeyValueTags(host.Tags).IgnoreConfig(ignoreTagsConfig).IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}
