package apigateway

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/apigateway"
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

// @SDKResource("aws_api_gateway_api_key")
func ResourceAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAPIKeyCreate,
		ReadWithoutTimeout:   resourceAPIKeyRead,
		UpdateWithoutTimeout: resourceAPIKeyUpdate,
		DeleteWithoutTimeout: resourceAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Default:  "Managed by Pulumi",
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"value": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(30, 128),
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceAPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(ctx, d.Get("tags").(map[string]interface{})))

	name := d.Get("name").(string)
	input := &apigateway.CreateApiKeyInput{
		Description: aws.String(d.Get("description").(string)),
		Enabled:     aws.Bool(d.Get("enabled").(bool)),
		Name:        aws.String(name),
		Tags:        Tags(tags.IgnoreAWS()),
		Value:       aws.String(d.Get("value").(string)),
	}

	apiKey, err := conn.CreateApiKeyWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating API Gateway API Key (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(apiKey.Id))

	return append(diags, resourceAPIKeyRead(ctx, d, meta)...)
}

func resourceAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	apiKey, err := FindAPIKeyByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] API Gateway API Key (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading API Gateway API Key (%s): %s", d.Id(), err)
	}

	tags := KeyValueTags(ctx, apiKey.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags_all: %s", err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "apigateway",
		Region:    meta.(*conns.AWSClient).Region,
		Resource:  fmt.Sprintf("/apikeys/%s", d.Id()),
	}.String()
	d.Set("arn", arn)

	d.Set("name", apiKey.Name)
	d.Set("description", apiKey.Description)
	d.Set("enabled", apiKey.Enabled)
	d.Set("value", apiKey.Value)

	if err := d.Set("created_date", apiKey.CreatedDate.Format(time.RFC3339)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting created_date: %s", err)
	}

	if err := d.Set("last_updated_date", apiKey.LastUpdatedDate.Format(time.RFC3339)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting last_updated_date: %s", err)
	}

	return diags
}

func resourceAPIKeyUpdateOperations(d *schema.ResourceData) []*apigateway.PatchOperation {
	operations := make([]*apigateway.PatchOperation, 0)
	if d.HasChange("enabled") {
		isEnabled := "false"
		if d.Get("enabled").(bool) {
			isEnabled = "true"
		}
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/enabled"),
			Value: aws.String(isEnabled),
		})
	}

	if d.HasChange("description") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/description"),
			Value: aws.String(d.Get("description").(string)),
		})
	}

	return operations
}

func resourceAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn()

	if d.HasChangesExcept("tags", "tags_all") {
		_, err := conn.UpdateApiKeyWithContext(ctx, &apigateway.UpdateApiKeyInput{
			ApiKey:          aws.String(d.Id()),
			PatchOperations: resourceAPIKeyUpdateOperations(d),
		})

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating API Gateway API Key (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(ctx, conn, d.Get("arn").(string), o, n); err != nil {
			return sdkdiag.AppendErrorf(diags, "updating tags: %s", err)
		}
	}

	return append(diags, resourceAPIKeyRead(ctx, d, meta)...)
}

func resourceAPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn()

	log.Printf("[DEBUG] Deleting API Gateway API Key: %s", d.Id())
	_, err := conn.DeleteApiKeyWithContext(ctx, &apigateway.DeleteApiKeyInput{
		ApiKey: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting API Gateway API Key (%s): %s", d.Id(), err)
	}

	return diags
}

func FindAPIKeyByID(ctx context.Context, conn *apigateway.APIGateway, id string) (*apigateway.ApiKey, error) {
	input := &apigateway.GetApiKeyInput{
		ApiKey:       aws.String(id),
		IncludeValue: aws.Bool(true),
	}

	output, err := conn.GetApiKeyWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
