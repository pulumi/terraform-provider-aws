package rds

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnapshotCreate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotUpdate,
		Delete: resourceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_snapshot_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"allocated_storage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_snapshot_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"license_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"option_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_db_snapshot_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))
	dBInstanceIdentifier := d.Get("db_instance_identifier").(string)

	params := &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(dBInstanceIdentifier),
		DBSnapshotIdentifier: aws.String(d.Get("db_snapshot_identifier").(string)),
		Tags:                 Tags(tags.IgnoreAWS()),
	}

	resp, err := conn.CreateDBSnapshot(params)
	if err != nil {
		return fmt.Errorf("Error creating AWS DB Snapshot %s: %s", dBInstanceIdentifier, err)
	}
	d.SetId(aws.StringValue(resp.DBSnapshot.DBSnapshotIdentifier))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    resourceSnapshotStateRefreshFunc(d, meta),
		Timeout:    d.Timeout(schema.TimeoutRead),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceSnapshotRead(d, meta)
}

func resourceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	params := &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: aws.String(d.Id()),
	}
	resp, err := conn.DescribeDBSnapshots(params)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBSnapshotNotFoundFault) {
		log.Printf("[WARN] AWS DB Snapshot (%s) is already gone", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error describing AWS DB Snapshot %s: %s", d.Id(), err)
	}

	snapshot := resp.DBSnapshots[0]

	arn := aws.StringValue(snapshot.DBSnapshotArn)
	d.Set("db_snapshot_identifier", snapshot.DBSnapshotIdentifier)
	d.Set("db_instance_identifier", snapshot.DBInstanceIdentifier)
	d.Set("allocated_storage", snapshot.AllocatedStorage)
	d.Set("availability_zone", snapshot.AvailabilityZone)
	d.Set("db_snapshot_arn", arn)
	d.Set("encrypted", snapshot.Encrypted)
	d.Set("engine", snapshot.Engine)
	d.Set("engine_version", snapshot.EngineVersion)
	d.Set("iops", snapshot.Iops)
	d.Set("kms_key_id", snapshot.KmsKeyId)
	d.Set("license_model", snapshot.LicenseModel)
	d.Set("option_group_name", snapshot.OptionGroupName)
	d.Set("port", snapshot.Port)
	d.Set("source_db_snapshot_identifier", snapshot.SourceDBSnapshotIdentifier)
	d.Set("source_region", snapshot.SourceRegion)
	d.Set("snapshot_type", snapshot.SnapshotType)
	d.Set("status", snapshot.Status)
	d.Set("vpc_id", snapshot.VpcId)

	tags, err := ListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("error listing tags for RDS DB Snapshot (%s): %s", arn, err)
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

func resourceSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	log.Printf("[DEBUG] Deleting RDS DB Snapshot: %s", d.Id())
	_, err := conn.DeleteDBSnapshot(&rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBSnapshotNotFoundFault) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting RDS DB Snapshot (%s): %w", d.Id(), err)
	}

	return nil
}

func resourceSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("db_snapshot_arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating RDS DB Snapshot (%s) tags: %s", d.Get("db_snapshot_arn").(string), err)
		}
	}

	return nil
}

func resourceSnapshotStateRefreshFunc(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn := meta.(*conns.AWSClient).RDSConn

		opts := &rds.DescribeDBSnapshotsInput{
			DBSnapshotIdentifier: aws.String(d.Id()),
		}

		log.Printf("[DEBUG] DB Snapshot describe configuration: %#v", opts)

		resp, err := conn.DescribeDBSnapshots(opts)
		if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBSnapshotNotFoundFault) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving DB Snapshots: %s", err)
		}

		if len(resp.DBSnapshots) != 1 {
			return nil, "", fmt.Errorf("No snapshots returned for %s", d.Id())
		}

		snapshot := resp.DBSnapshots[0]

		return resp, *snapshot.Status, nil
	}
}
