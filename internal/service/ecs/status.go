package ecs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	serviceStatusInactive = "INACTIVE"
	serviceStatusActive   = "ACTIVE"
	serviceStatusDraining = "DRAINING"

	serviceStatusError = "ERROR"
	serviceStatusNone  = "NONE"

	clusterStatusError = "ERROR"
	clusterStatusNone  = "NONE"

	taskSetStatusActive   = "ACTIVE"
	taskSetStatusDraining = "DRAINING"
	taskSetStatusPrimary  = "PRIMARY"
)

func statusCapacityProvider(conn *ecs.ECS, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindCapacityProviderByARN(conn, arn)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusCapacityProviderUpdate(conn *ecs.ECS, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindCapacityProviderByARN(conn, arn)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.UpdateStatus), nil
	}
}

func statusServiceNoTags(conn *ecs.ECS, id, cluster string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := FindServiceNoTagsByID(context.TODO(), conn, id, cluster)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", err
		}

		return service, aws.StringValue(service.Status), err
	}
}

func statusService(conn *ecs.ECS, id, cluster string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &ecs.DescribeServicesInput{
			Services: aws.StringSlice([]string{id}),
			Cluster:  aws.String(cluster),
		}

		output, err := conn.DescribeServices(input)

		if tfawserr.ErrCodeEquals(err, ecs.ErrCodeServiceNotFoundException) {
			return nil, serviceStatusNone, nil
		}

		if err != nil {
			return nil, serviceStatusError, err
		}

		if len(output.Services) == 0 {
			return nil, serviceStatusNone, nil
		}

		log.Printf("[DEBUG] ECS service (%s) is currently %q", id, *output.Services[0].Status)
		return output, aws.StringValue(output.Services[0].Status), err
	}
}

func statusCluster(ctx context.Context, conn *ecs.ECS, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := FindClusterByNameOrARN(ctx, conn, arn)

		if tfresource.NotFound(err) {
			return nil, clusterStatusNone, nil
		}

		if err != nil {
			return nil, clusterStatusError, err
		}

		return cluster, aws.StringValue(cluster.Status), err
	}
}

func stabilityStatusTaskSet(conn *ecs.ECS, taskSetID, service, cluster string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &ecs.DescribeTaskSetsInput{
			Cluster:  aws.String(cluster),
			Service:  aws.String(service),
			TaskSets: aws.StringSlice([]string{taskSetID}),
		}

		output, err := conn.DescribeTaskSets(input)

		if err != nil {
			return nil, "", err
		}

		if output == nil || len(output.TaskSets) == 0 {
			return nil, "", nil
		}

		return output.TaskSets[0], aws.StringValue(output.TaskSets[0].StabilityStatus), nil
	}
}

func statusTaskSet(conn *ecs.ECS, taskSetID, service, cluster string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &ecs.DescribeTaskSetsInput{
			Cluster:  aws.String(cluster),
			Service:  aws.String(service),
			TaskSets: aws.StringSlice([]string{taskSetID}),
		}

		output, err := conn.DescribeTaskSets(input)

		if err != nil {
			return nil, "", err
		}

		if output == nil || len(output.TaskSets) == 0 {
			return nil, "", nil
		}

		return output.TaskSets[0], aws.StringValue(output.TaskSets[0].Status), nil
	}
}
