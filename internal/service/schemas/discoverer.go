package schemas

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/schemas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceDiscoverer() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscovererCreate,
		Read:   resourceDiscovererRead,
		Update: resourceDiscovererUpdate,
		Delete: resourceDiscovererDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},

			"source_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceDiscovererCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	sourceARN := d.Get("source_arn").(string)
	input := &schemas.CreateDiscovererInput{
		SourceArn: aws.String(sourceARN),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	log.Printf("[DEBUG] Creating EventBridge Schemas Discoverer: %s", input)
	output, err := conn.CreateDiscoverer(input)

	if err != nil {
		return fmt.Errorf("error creating EventBridge Schemas Discoverer (%s): %w", sourceARN, err)
	}

	d.SetId(aws.StringValue(output.DiscovererId))

	return resourceDiscovererRead(d, meta)
}

func resourceDiscovererRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	output, err := FindDiscovererByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EventBridge Schemas Discoverer (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EventBridge Schemas Discoverer (%s): %w", d.Id(), err)
	}

	d.Set("arn", output.DiscovererArn)
	d.Set("description", output.Description)
	d.Set("source_arn", output.SourceArn)

	tags, err := ListTags(conn, d.Get("arn").(string))

	if err != nil {
		return fmt.Errorf("error listing tags for EventBridge Schemas Discoverer (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceDiscovererUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn

	if d.HasChange("description") {
		input := &schemas.UpdateDiscovererInput{
			DiscovererId: aws.String(d.Id()),
			Description:  aws.String(d.Get("description").(string)),
		}

		log.Printf("[DEBUG] Updating EventBridge Schemas Discoverer: %s", input)
		_, err := conn.UpdateDiscoverer(input)

		if err != nil {
			return fmt.Errorf("error updating EventBridge Schemas Discoverer (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %w", err)
		}
	}

	return resourceDiscovererRead(d, meta)
}

func resourceDiscovererDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn

	log.Printf("[INFO] Deleting EventBridge Schemas Discoverer (%s)", d.Id())
	_, err := conn.DeleteDiscoverer(&schemas.DeleteDiscovererInput{
		DiscovererId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, schemas.ErrCodeNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting EventBridge Schemas Discoverer (%s): %w", d.Id(), err)
	}

	return nil
}
