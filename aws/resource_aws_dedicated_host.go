package aws

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsDedicatedHostCreate,
		ReadContext:   resourceAwsDedicatedHostRead,
		UpdateContext: resourceAwsDedicatedHostUpdate,
		DeleteContext: resourceAwsDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		SchemaVersion: 1,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"tags": tagsSchema(),
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"instance_type", "instance_family"},
			},
			"instance_family": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"instance_type", "instance_family"},
			},
			"host_recovery": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					ec2.HostRecoveryOn,
					ec2.HostRecoveryOff,
				}, false),
			},
			"auto_placement": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					ec2.AutoPlacementOn,
					ec2.AutoPlacementOff,
				}, false),
			},
		},
	}
}

// resourceAwsDedicatedHostCreate allocates a Dedicated Host to your account.
// https://docs.aws.amazon.com/en_pv/AWSEC2/latest/APIReference/API_AllocateHosts.html
func resourceAwsDedicatedHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).ec2conn

	runOpts := &ec2.AllocateHostsInput{
		AvailabilityZone: aws.String(d.Get("availability_zone").(string)),
		Quantity:         aws.Int64(int64(1)),
	}
	if v, ok := d.GetOk("instance_type"); ok {
		runOpts.InstanceType = aws.String(v.(string))
	}
	if v, ok := d.GetOk("instance_family"); ok {
		runOpts.InstanceFamily = aws.String(v.(string))
	}
	if v, ok := d.GetOk("host_recovery"); ok {
		runOpts.HostRecovery = aws.String(v.(string))
	}
	if v, ok := d.GetOk("auto_placement"); ok {
		runOpts.AutoPlacement = aws.String(v.(string))
	}

	tagsSpec := ec2TagSpecificationsFromMap(d.Get("tags").(map[string]interface{}), ec2.ResourceTypeDedicatedHost)
	if len(tagsSpec) > 0 {
		runOpts.TagSpecifications = tagsSpec
	}

	runResp, err := conn.AllocateHostsWithContext(ctx, runOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	if runResp == nil || len(runResp.HostIds) == 0 {
		return diag.Errorf("Error launching source host: no hosts returned in response")
	}

	log.Printf("[INFO] Host ID: %s", *runResp.HostIds[0])
	d.SetId(*runResp.HostIds[0])

	return resourceAwsDedicatedHostRead(ctx, d, meta)
}

func resourceAwsDedicatedHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).ec2conn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	resp, err := conn.DescribeHostsWithContext(ctx, &ec2.DescribeHostsInput{
		HostIds: []*string{aws.String(d.Id())},
	})
	if err != nil {
		// If the host was not found, return nil so that we can show
		// that host is gone.
		if isAWSErr(err, "InvalidHostID.NotFound", "") {
			d.SetId("")
			return nil
		}

		// Some other error, report it
		return diag.FromErr(err)
	}
	if len(resp.Hosts) == 0 {
		d.SetId("")
		return nil
	}

	host := resp.Hosts[0]
	if host.HostProperties == nil {
		return diag.Errorf("HostProperties for dedicated host is nil")
	}

	d.Set("auto_placement", host.AutoPlacement)
	d.Set("availability_zone", host.AvailabilityZone)
	d.Set("host_recovery", host.HostRecovery)
	d.Set("instance_type", host.HostProperties.InstanceType)
	d.Set("instance_family", host.HostProperties.InstanceFamily)

	if err := d.Set("tags", keyvaluetags.Ec2KeyValueTags(host.Tags).IgnoreConfig(ignoreTagsConfig).IgnoreAws().Map()); err != nil {
		return diag.Errorf("error setting tags: %s", err)
	}

	return nil
}

// resourceAwsDedicatedHostUpdate modifies AWS Host AutoPlacement and HostRecovery settings.
// When auto-placement is enabled, any instances that you launch with a tenancy of host but without a specific host ID are placed onto any available
// Dedicated Host in your account that has auto-placement enabled.
// https://docs.aws.amazon.com/en_pv/AWSEC2/latest/APIReference/API_ModifyHosts.html
func resourceAwsDedicatedHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).ec2conn
	requestUpdate := false

	req := &ec2.ModifyHostsInput{
		HostIds: []*string{aws.String(d.Id())},
	}

	if d.HasChange("auto_placement") {
		req.AutoPlacement = aws.String(d.Get("auto_placement").(string))
		requestUpdate = true
	}

	if d.HasChange("host_recovery") {
		req.HostRecovery = aws.String(d.Get("host_recovery").(string))
		requestUpdate = true
	}

	if d.HasChange("instance_family") {
		instanceFamily, ok := d.GetOk("instance_family")
		if ok {
			req.InstanceFamily = aws.String(instanceFamily.(string))
		}
		requestUpdate = true
	}
	if d.HasChange("instance_type") {
		instanceType, ok := d.GetOk("instance_type")
		if ok {
			req.InstanceType = aws.String(instanceType.(string))
		}
		requestUpdate = true
	}

	if requestUpdate {
		_, err := conn.ModifyHostsWithContext(ctx, req)
		if err != nil {
			return diag.Errorf("Error modifying host %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")

		if err := keyvaluetags.Ec2UpdateTags(conn, d.Id(), o, n); err != nil {
			return diag.Errorf("error updating tags: %s", err)
		}
	}

	return resourceAwsDedicatedHostRead(ctx, d, meta)
}

func resourceAwsDedicatedHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).ec2conn

	if _, err := conn.ReleaseHostsWithContext(ctx, &ec2.ReleaseHostsInput{
		HostIds: []*string{aws.String(d.Id())},
	}); err != nil {
		return diag.Errorf("error terminating EC2 Host %q: %s", d.Id(), err)
	}

	return nil
}
