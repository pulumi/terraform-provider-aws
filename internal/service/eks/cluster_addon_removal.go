package eks

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func removeAddons(d *schema.ResourceData, conn *eks.EKS) error {
	if v, ok := d.GetOk("default_addons_to_remove"); ok && len(v.([]interface{})) > 0 {
		ctx := context.Background()
		var wg sync.WaitGroup
		var removalErrors *multierror.Error

		for _, addon := range flex.ExpandStringList(v.([]interface{})) {
			if addon == nil {
				return fmt.Errorf("addonName cannot be dereferenced")
			}
			addonName := *addon
			wg.Add(1)

			go func() {
				defer wg.Done()
				removalErrors = multierror.Append(removalErrors, removeAddon(d, conn, addonName, ctx))
			}()
		}
		wg.Wait()
		return removalErrors.ErrorOrNil()
	}
	return nil
}

func removeAddon(d *schema.ResourceData, conn *eks.EKS, addonName string, ctx context.Context) error {
	log.Printf("[DEBUG] Creating EKS Add-On: %s", addonName)
	createAddonInput := &eks.CreateAddonInput{
		AddonName:          aws.String(addonName),
		ClientRequestToken: aws.String(resource.UniqueId()),
		ClusterName:        aws.String(d.Id()),
		ResolveConflicts:   aws.String(eks.ResolveConflictsOverwrite),
	}

	err := resource.RetryContext(ctx, propagationTimeout, func() *resource.RetryError {
		_, err := conn.CreateAddonWithContext(ctx, createAddonInput)

		if tfawserr.ErrMessageContains(err, eks.ErrCodeInvalidParameterException, "CREATE_FAILED") {
			return resource.RetryableError(err)
		}

		if tfawserr.ErrMessageContains(err, eks.ErrCodeInvalidParameterException, "does not exist") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.CreateAddonWithContext(ctx, createAddonInput)
	}

	if err != nil {
		return fmt.Errorf("error creating EKS Add-On (%s): %w", addonName, err)
	}

	_, err = waitAddonCreatedAllowDegraded(ctx, conn, d.Id(), addonName)

	if err != nil {
		return fmt.Errorf("unexpected EKS Add-On (%s) state returned during creation: %w", addonName, err)
	}
	log.Printf("[DEBUG] Created EKS Add-On: %s", addonName)

	deleteAddonInput := &eks.DeleteAddonInput{
		AddonName:   aws.String(addonName),
		ClusterName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting EKS Add-On: %s", addonName)
	_, err = conn.DeleteAddonWithContext(ctx, deleteAddonInput)

	if err != nil {
		return fmt.Errorf("error deleting EKS Add-On (%s): %w", addonName, err)
	}

	_, err = waitAddonDeleted(ctx, conn, d.Id(), addonName, addonDeletedTimeout)

	if err != nil {
		return fmt.Errorf("error waiting for EKS Add-On (%s) to delete: %w", addonName, err)
	}
	log.Printf("[DEBUG] Deleted EKS Add-On: %s", addonName)
	return nil
}

func waitAddonCreatedAllowDegraded(ctx context.Context, conn *eks.EKS, clusterName, addonName string) (*eks.Addon, error) {
	// We do not care about the addons actually being created successfully here. We only want them to be adopted by
	// Terraform to be able to fully remove them afterwards again.
	stateConf := resource.StateChangeConf{
		Pending: []string{eks.AddonStatusCreating},
		Target:  []string{eks.AddonStatusActive, eks.AddonStatusDegraded},
		Refresh: statusAddon(ctx, conn, clusterName, addonName),
		Timeout: addonCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.Addon); ok {
		if status, health := aws.StringValue(output.Status), output.Health; status == eks.AddonStatusCreateFailed && health != nil {
			tfresource.SetLastError(err, AddonIssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}
