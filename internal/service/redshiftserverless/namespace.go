package redshiftserverless

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceNamespaceCreate,
		Read:   resourceNamespaceRead,
		Update: resourceNamespaceUpdate,
		Delete: resourceNamespaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"admin_user_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"admin_username": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"default_iam_role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},
			"iam_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidARN,
			},
			"log_exports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"userlog", "connectionlog", "useractivitylog"}, false),
				},
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("namespace_name").(string)
	input := &redshiftserverless.CreateNamespaceInput{
		NamespaceName: aws.String(name),
		Tags:          Tags(tags.IgnoreAWS()),
	}

	if v, ok := d.GetOk("admin_user_password"); ok {
		input.AdminUserPassword = aws.String(v.(string))
	}

	if v, ok := d.GetOk("admin_username"); ok {
		input.AdminUsername = aws.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		input.DbName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("default_iam_role_arn"); ok {
		input.DefaultIamRoleArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("iam_roles"); ok && v.(*schema.Set).Len() > 0 {
		input.IamRoles = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		input.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("log_exports"); ok && v.(*schema.Set).Len() > 0 {
		input.LogExports = flex.ExpandStringSet(v.(*schema.Set))
	}

	output, err := conn.CreateNamespace(input)

	if err != nil {
		return fmt.Errorf("creating Redshift Serverless Namespace (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.Namespace.NamespaceName))

	return resourceNamespaceRead(d, meta)
}

func resourceNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn()
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	output, err := FindNamespaceByName(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Redshift Serverless Namespace (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading Redshift Serverless Namespace (%s): %w", d.Id(), err)
	}

	arn := aws.StringValue(output.NamespaceArn)
	d.Set("admin_username", output.AdminUsername)
	d.Set("arn", arn)
	d.Set("db_name", output.DbName)
	d.Set("default_iam_role_arn", output.DefaultIamRoleArn)
	d.Set("iam_roles", aws.StringValueSlice(output.IamRoles))
	d.Set("kms_key_id", output.KmsKeyId)
	d.Set("log_exports", aws.StringValueSlice(output.LogExports))
	d.Set("namespace_id", output.NamespaceId)
	d.Set("namespace_name", output.NamespaceName)

	tags, err := ListTags(conn, arn)

	if tfawserr.ErrCodeEquals(err, "UnknownOperationException") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("listing tags for Redshift Serverless Namespace (%s): %w", arn, err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	return nil
}

func resourceNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn()

	if d.HasChangesExcept("tags", "tags_all") {
		input := &redshiftserverless.UpdateNamespaceInput{
			NamespaceName: aws.String(d.Id()),
		}

		if d.HasChanges("admin_username", "admin_user_password") {
			input.AdminUsername = aws.String(d.Get("admin_username").(string))
			input.AdminUserPassword = aws.String(d.Get("admin_user_password").(string))
		}

		if d.HasChange("default_iam_role_arn") {
			input.DefaultIamRoleArn = aws.String(d.Get("default_iam_role_arn").(string))
		}

		if d.HasChange("iam_roles") {
			input.IamRoles = flex.ExpandStringSet(d.Get("iam_roles").(*schema.Set))
		}

		if d.HasChange("kms_key_id") {
			input.KmsKeyId = aws.String(d.Get("kms_key_id").(string))
		}

		if d.HasChange("log_exports") {
			input.LogExports = flex.ExpandStringSet(d.Get("log_exports").(*schema.Set))
		}

		_, err := conn.UpdateNamespace(input)

		if err != nil {
			return fmt.Errorf("updating Redshift Serverless Namespace (%s): %w", d.Id(), err)
		}

		if _, err := waitNamespaceUpdated(conn, d.Id()); err != nil {
			return fmt.Errorf("waiting for Redshift Serverless Namespace (%s) update: %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("updating Redshift Serverless Namespace (%s) tags: %w", d.Get("arn").(string), err)
		}
	}

	return resourceNamespaceRead(d, meta)
}

func resourceNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn()

	log.Printf("[DEBUG] Deleting Redshift Serverless Namespace: %s", d.Id())
	_, err := tfresource.RetryWhenAWSErrMessageContains(10*time.Minute,
		func() (interface{}, error) {
			return conn.DeleteNamespace(&redshiftserverless.DeleteNamespaceInput{
				NamespaceName: aws.String(d.Id()),
			})
		},
		// "ConflictException: There is an operation running on the namespace. Try deleting the namespace again later."
		redshiftserverless.ErrCodeConflictException, "operation running")

	if tfawserr.ErrCodeEquals(err, redshiftserverless.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting Redshift Serverless Namespace (%s): %w", d.Id(), err)
	}

	if _, err := waitNamespaceDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("waiting for Redshift Serverless Namespace (%s) delete: %w", d.Id(), err)
	}

	return nil
}
