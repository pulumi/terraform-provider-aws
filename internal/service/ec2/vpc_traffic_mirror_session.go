package ec2

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceTrafficMirrorSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficMirrorSessionCreate,
		Update: resourceTrafficMirrorSessionUpdate,
		Read:   resourceTrafficMirrorSessionRead,
		Delete: resourceTrafficMirrorSessionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"packet_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"session_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 32766),
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_target_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_network_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16777216),
			},
		},
	}
}

func resourceTrafficMirrorSessionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &ec2.CreateTrafficMirrorSessionInput{
		NetworkInterfaceId:    aws.String(d.Get("network_interface_id").(string)),
		TrafficMirrorFilterId: aws.String(d.Get("traffic_mirror_filter_id").(string)),
		TrafficMirrorTargetId: aws.String(d.Get("traffic_mirror_target_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("packet_length"); ok {
		input.PacketLength = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("session_number"); ok {
		input.SessionNumber = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("virtual_network_id"); ok {
		input.VirtualNetworkId = aws.Int64(int64(v.(int)))
	}

	if len(tags) > 0 {
		input.TagSpecifications = tagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeTrafficMirrorSession)
	}

	output, err := conn.CreateTrafficMirrorSession(input)

	if err != nil {
		return fmt.Errorf("creating EC2 Traffic Mirror Session: %w", err)
	}

	d.SetId(aws.StringValue(output.TrafficMirrorSession.TrafficMirrorSessionId))

	return resourceTrafficMirrorSessionRead(d, meta)
}

func resourceTrafficMirrorSessionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	session, err := FindTrafficMirrorSessionByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Traffic Mirror Session %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading EC2 Traffic Mirror Session (%s): %w", d.Id(), err)
	}

	ownerID := aws.StringValue(session.OwnerId)
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: ownerID,
		Resource:  fmt.Sprintf("traffic-mirror-session/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("description", session.Description)
	d.Set("network_interface_id", session.NetworkInterfaceId)
	d.Set("owner_id", ownerID)
	d.Set("packet_length", session.PacketLength)
	d.Set("session_number", session.SessionNumber)
	d.Set("traffic_mirror_filter_id", session.TrafficMirrorFilterId)
	d.Set("traffic_mirror_target_id", session.TrafficMirrorTargetId)
	d.Set("virtual_network_id", session.VirtualNetworkId)

	tags := KeyValueTags(session.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	return nil
}

func resourceTrafficMirrorSessionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	if d.HasChangesExcept("tags", "tags_all") {
		input := &ec2.ModifyTrafficMirrorSessionInput{
			TrafficMirrorSessionId: aws.String(d.Id()),
		}

		if d.HasChange("session_number") {
			input.SessionNumber = aws.Int64(int64(d.Get("session_number").(int)))
		}

		if d.HasChange("traffic_mirror_filter_id") {
			input.TrafficMirrorFilterId = aws.String(d.Get("traffic_mirror_filter_id").(string))
		}

		if d.HasChange("traffic_mirror_target_id") {
			input.TrafficMirrorTargetId = aws.String(d.Get("traffic_mirror_target_id").(string))
		}

		var removeFields []string

		if d.HasChange("description") {
			if v := d.Get("description").(string); v != "" {
				input.Description = aws.String(v)
			} else {
				removeFields = append(removeFields, ec2.TrafficMirrorSessionFieldDescription)
			}
		}

		if d.HasChange("packet_length") {
			if v := d.Get("packet_length").(int); v != 0 {
				input.PacketLength = aws.Int64(int64(v))
			} else {
				removeFields = append(removeFields, ec2.TrafficMirrorSessionFieldPacketLength)
			}
		}

		if d.HasChange("virtual_network_id") {
			if v := d.Get("virtual_network_id").(int); v != 0 {
				input.VirtualNetworkId = aws.Int64(int64(v))
			} else {
				removeFields = append(removeFields, ec2.TrafficMirrorSessionFieldVirtualNetworkId)
			}
		}

		if len(removeFields) > 0 {
			input.RemoveFields = aws.StringSlice(removeFields)
		}

		_, err := conn.ModifyTrafficMirrorSession(input)

		if err != nil {
			return fmt.Errorf("updating EC2 Traffic Mirror Session (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("updating EC2 Traffic Mirror Session (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceTrafficMirrorSessionRead(d, meta)
}

func resourceTrafficMirrorSessionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	log.Printf("[DEBUG] Deleting EC2 Traffic Mirror Session: %s", d.Id())
	_, err := conn.DeleteTrafficMirrorSession(&ec2.DeleteTrafficMirrorSessionInput{
		TrafficMirrorSessionId: aws.String(d.Id()),
	})

	if err != nil {
		return fmt.Errorf("deleting EC2 Traffic Mirror Session (%s): %w", d.Id(), err)
	}

	return nil
}
