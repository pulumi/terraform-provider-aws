package datapipeline

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/datapipeline"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DataPipelineConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	uniqueID := resource.UniqueId()

	input := datapipeline.CreatePipelineInput{
		Name:     aws.String(d.Get("name").(string)),
		UniqueId: aws.String(uniqueID),
		Tags:     Tags(tags.IgnoreAWS()),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	resp, err := conn.CreatePipeline(&input)

	if err != nil {
		return fmt.Errorf("Error creating datapipeline: %s", err)
	}

	d.SetId(aws.StringValue(resp.PipelineId))

	return resourcePipelineRead(d, meta)
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DataPipelineConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	v, err := PipelineRetrieve(d.Id(), conn)
	if tfawserr.ErrCodeEquals(err, datapipeline.ErrCodePipelineNotFoundException) || tfawserr.ErrCodeEquals(err, datapipeline.ErrCodePipelineDeletedException) || v == nil {
		log.Printf("[WARN] DataPipeline (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error describing DataPipeline (%s): %s", d.Id(), err)
	}

	d.Set("name", v.Name)
	d.Set("description", v.Description)
	tags := KeyValueTags(v.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DataPipelineConn()

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating Datapipeline Pipeline (%s) tags: %s", d.Id(), err)
		}
	}

	return resourcePipelineRead(d, meta)
}

func resourcePipelineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DataPipelineConn()

	opts := datapipeline.DeletePipelineInput{
		PipelineId: aws.String(d.Id()),
	}

	_, err := conn.DeletePipeline(&opts)
	if tfawserr.ErrCodeEquals(err, datapipeline.ErrCodePipelineNotFoundException) || tfawserr.ErrCodeEquals(err, datapipeline.ErrCodePipelineDeletedException) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error deleting Data Pipeline %s: %s", d.Id(), err.Error())
	}

	return WaitForDeletion(conn, d.Id())
}

func PipelineRetrieve(id string, conn *datapipeline.DataPipeline) (*datapipeline.PipelineDescription, error) {
	opts := datapipeline.DescribePipelinesInput{
		PipelineIds: []*string{aws.String(id)},
	}

	resp, err := conn.DescribePipelines(&opts)
	if err != nil {
		return nil, err
	}

	var pipeline *datapipeline.PipelineDescription

	for _, p := range resp.PipelineDescriptionList {
		if p == nil {
			continue
		}

		if aws.StringValue(p.PipelineId) == id {
			pipeline = p
			break
		}
	}

	return pipeline, nil
}

func WaitForDeletion(conn *datapipeline.DataPipeline, pipelineID string) error {
	params := &datapipeline.DescribePipelinesInput{
		PipelineIds: []*string{aws.String(pipelineID)},
	}
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := conn.DescribePipelines(params)
		if tfawserr.ErrCodeEquals(err, datapipeline.ErrCodePipelineNotFoundException) || tfawserr.ErrCodeEquals(err, datapipeline.ErrCodePipelineDeletedException) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("DataPipeline (%s) still exists", pipelineID))
	})
}
