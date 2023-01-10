package schemas

import (
	"fmt"
	"log"
	"regexp"
	"time"

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

func ResourceSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchemaCreate,
		Read:   resourceSchemaRead,
		Update: resourceSchemaUpdate,
		Delete: resourceSchemaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"content": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},

			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 385),
					validation.StringMatch(regexp.MustCompile(`^[\.\-_A-Za-z@]+`), ""),
				),
			},

			"registry_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(schemas.Type_Values(), true),
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version_created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceSchemaCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("name").(string)
	registryName := d.Get("registry_name").(string)
	input := &schemas.CreateSchemaInput{
		Content:      aws.String(d.Get("content").(string)),
		RegistryName: aws.String(registryName),
		SchemaName:   aws.String(name),
		Type:         aws.String(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	id := SchemaCreateResourceID(name, registryName)

	log.Printf("[DEBUG] Creating EventBridge Schemas Schema: %s", input)
	_, err := conn.CreateSchema(input)

	if err != nil {
		return fmt.Errorf("error creating EventBridge Schemas Schema (%s): %w", id, err)
	}

	d.SetId(id)

	return resourceSchemaRead(d, meta)
}

func resourceSchemaRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name, registryName, err := SchemaParseResourceID(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing EventBridge Schemas Schema ID: %w", err)
	}

	output, err := FindSchemaByNameAndRegistryName(conn, name, registryName)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EventBridge Schemas Schema (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EventBridge Schemas Schema (%s): %w", d.Id(), err)
	}

	d.Set("arn", output.SchemaArn)
	d.Set("content", output.Content)
	d.Set("description", output.Description)
	if output.LastModified != nil {
		d.Set("last_modified", aws.TimeValue(output.LastModified).Format(time.RFC3339))
	} else {
		d.Set("last_modified", nil)
	}
	d.Set("name", output.SchemaName)
	d.Set("registry_name", registryName)
	d.Set("type", output.Type)
	d.Set("version", output.SchemaVersion)
	if output.VersionCreatedDate != nil {
		d.Set("version_created_date", aws.TimeValue(output.VersionCreatedDate).Format(time.RFC3339))
	} else {
		d.Set("version_created_date", nil)
	}

	tags, err := ListTags(conn, d.Get("arn").(string))

	if err != nil {
		return fmt.Errorf("error listing tags for EventBridge Schemas Schema (%s): %w", d.Id(), err)
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

func resourceSchemaUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn()

	if d.HasChanges("content", "description", "type") {
		name, registryName, err := SchemaParseResourceID(d.Id())

		if err != nil {
			return fmt.Errorf("error parsing EventBridge Schemas Schema ID: %w", err)
		}

		input := &schemas.UpdateSchemaInput{
			RegistryName: aws.String(registryName),
			SchemaName:   aws.String(name),
		}

		if d.HasChanges("content", "type") {
			input.Content = aws.String(d.Get("content").(string))
			input.Type = aws.String(d.Get("type").(string))
		}

		if d.HasChange("description") {
			input.Description = aws.String(d.Get("description").(string))
		}

		log.Printf("[DEBUG] Updating EventBridge Schemas Schema: %s", input)
		_, err = conn.UpdateSchema(input)

		if err != nil {
			return fmt.Errorf("error updating EventBridge Schemas Schema (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %w", err)
		}
	}

	return resourceSchemaRead(d, meta)
}

func resourceSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SchemasConn()

	name, registryName, err := SchemaParseResourceID(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing EventBridge Schemas Schema ID: %w", err)
	}

	log.Printf("[INFO] Deleting EventBridge Schemas Schema (%s)", d.Id())
	_, err = conn.DeleteSchema(&schemas.DeleteSchemaInput{
		RegistryName: aws.String(registryName),
		SchemaName:   aws.String(name),
	})

	if tfawserr.ErrCodeEquals(err, schemas.ErrCodeNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting EventBridge Schemas Schema (%s): %w", d.Id(), err)
	}

	return nil
}
