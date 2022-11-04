package ecr

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-aws/internal/conns"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCredentials() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsEcrCredentialsRead,

		Schema: map[string]*schema.Schema{
			"registry_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"authorization_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsEcrCredentialsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ECRConn()

	registryID := d.Get("registry_id").(string)
	log.Printf("[DEBUG] Reading ECR repository credentials %s", registryID)

	out, err := conn.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(registryID)},
	})

	if err != nil {
		if ecrerr, ok := err.(awserr.Error); ok && ecrerr.Code() == "ErrCodeInvalidParameterException" {
			log.Printf("[WARN] ECR Repository %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	auth := out.AuthorizationData[0]
	log.Printf("[DEBUG] Received ECR repository credentials for %s", auth.ProxyEndpoint)

	d.SetId(registryID)
	d.Set("authorization_token", auth.AuthorizationToken)
	d.Set("expires_at", auth.ExpiresAt.Format(time.RFC3339))
	d.Set("proxy_endpoint", auth.ProxyEndpoint)

	return nil
}
