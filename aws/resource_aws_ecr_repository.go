package aws

import (
	"log"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsEcrRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsEcrRepositoryCreate,
		Read:   resourceAwsEcrRepositoryRead,
		Delete: resourceAwsEcrRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsEcrRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ecrconn

	input := ecr.CreateRepositoryInput{
		RepositoryName: aws.String(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Creating ECR resository: %s", input)
	out, err := conn.CreateRepository(&input)
	if err != nil {
		return err
	}

	repository := *out.Repository

	log.Printf("[DEBUG] ECR repository created: %q", *repository.RepositoryArn)

	d.SetId(*repository.RepositoryName)
	d.Set("arn", repository.RepositoryArn)
	d.Set("registry_id", repository.RegistryId)

	// Ensure AWS will give us back the registry info before calling read.
	log.Printf("[DEBUG] Waiting for ECR repository (%s) to exist", *repository.RepositoryArn)
	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"exists"},
		Refresh: resourceAwsEcrRepositoryRefreshFunc(conn, d.Id()),
		Timeout: 3 * time.Minute,
	}
	repRaw, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for ECR repository (%s) to become available: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] ECR repository (%s) exists", *repository.RepositoryArn)
	resourceAwsEcrRepositoryReadData(d, repRaw.(*ecr.Repository), meta)
	return nil
}

func resourceAwsEcrRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading repository %s", d.Id())

	conn := meta.(*AWSClient).ecrconn
	repRaw, _, err := resourceAwsEcrRepositoryRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if repRaw == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Received repository %s", repRaw)
	resourceAwsEcrRepositoryReadData(d, repRaw.(*ecr.Repository), meta)
	return nil
}

func resourceAwsEcrRepositoryReadData(d *schema.ResourceData, repository *ecr.Repository, meta interface{}) {
	d.SetId(*repository.RepositoryName)
	d.Set("arn", repository.RepositoryArn)
	d.Set("registry_id", repository.RegistryId)
	d.Set("name", repository.RepositoryName)

	repositoryUrl := buildRepositoryUrl(repository, meta.(*AWSClient).region)
	log.Printf("[INFO] Setting the repository url to be %s", repositoryUrl)
	d.Set("repository_url", repositoryUrl)
}

func buildRepositoryUrl(repo *ecr.Repository, region string) string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s", *repo.RegistryId, region, *repo.RepositoryName)
}

func resourceAwsEcrRepositoryRefreshFunc(conn *ecr.ECR, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeRepositories(&ecr.DescribeRepositoriesInput{
			RepositoryNames: []*string{aws.String(id)},
		})
		if err != nil {
			if ecrerr, ok := err.(awserr.Error); ok && ecrerr.Code() == "RepositoryNotFoundException" {
				resp = nil
				err = nil
			} else {
				return nil, "", err
			}
		}
		if resp == nil {
			return nil, "", nil
		}

		return resp.Repositories[0], "exists", nil
	}
}

func resourceAwsEcrRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ecrconn

	_, err := conn.DeleteRepository(&ecr.DeleteRepositoryInput{
		RepositoryName: aws.String(d.Id()),
		RegistryId:     aws.String(d.Get("registry_id").(string)),
		Force:          aws.Bool(true),
	})
	if err != nil {
		if ecrerr, ok := err.(awserr.Error); ok && ecrerr.Code() == "RepositoryNotFoundException" {
			return nil
		}
		return err
	}

	log.Printf("[DEBUG] Waiting for ECR Repository %q to be deleted", d.Id())
	err = resource.Retry(20*time.Minute, func() *resource.RetryError {
		_, err := conn.DescribeRepositories(&ecr.DescribeRepositoriesInput{
			RepositoryNames: []*string{aws.String(d.Id())},
		})

		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if !ok {
				return resource.NonRetryableError(err)
			}

			if awsErr.Code() == "RepositoryNotFoundException" {
				return nil
			}

			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(
			fmt.Errorf("%q: Timeout while waiting for the ECR Repository to be deleted", d.Id()))
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] repository %q deleted.", d.Get("name").(string))

	return nil
}
