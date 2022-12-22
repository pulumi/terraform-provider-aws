package simpledb

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/simpledb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Delete: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SimpleDBConn

	name := d.Get("name").(string)
	input := &simpledb.CreateDomainInput{
		DomainName: aws.String(name),
	}
	_, err := conn.CreateDomain(input)
	if err != nil {
		return fmt.Errorf("Create SimpleDB Domain failed: %s", err)
	}

	d.SetId(name)
	return nil
}

func resourceDomainRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SimpleDBConn

	input := &simpledb.DomainMetadataInput{
		DomainName: aws.String(d.Id()),
	}
	_, err := conn.DomainMetadata(input)
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == "NoSuchDomain" {
			log.Printf("[WARN] Removing SimpleDB domain %q because it's gone.", d.Id())
			d.SetId("")
			return nil
		}
	}
	if err != nil {
		return err
	}

	d.Set("name", d.Id())
	return nil
}

func resourceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SimpleDBConn

	input := &simpledb.DeleteDomainInput{
		DomainName: aws.String(d.Id()),
	}
	_, err := conn.DeleteDomain(input)
	if err != nil {
		return fmt.Errorf("Delete SimpleDB Domain failed: %s", err)
	}

	return nil
}
