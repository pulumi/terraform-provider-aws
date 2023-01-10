---
subcategory: "SageMaker"
layout: "aws"
page_title: "AWS: aws_sagemaker_user_profile"
description: |-
  Provides a SageMaker User Profile resource.
---

# Resource: aws_sagemaker_user_profile

Provides a SageMaker User Profile resource.

## Example Usage

### Basic usage

```terraform
resource "aws_sagemaker_user_profile" "example" {
  domain_id         = aws_sagemaker_domain.test.id
  user_profile_name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `user_profile_name` - (Required) The name for the User Profile.
* `domain_id` - (Required) The ID of the associated Domain.
* `single_sign_on_user_identifier` - (Optional) A specifier for the type of value specified in `single_sign_on_user_value`. Currently, the only supported value is `UserName`. If the Domain's AuthMode is SSO, this field is required. If the Domain's AuthMode is not SSO, this field cannot be specified.
* `single_sign_on_user_value` - (Required) The username of the associated AWS Single Sign-On User for this User Profile. If the Domain's AuthMode is SSO, this field is required, and must match a valid username of a user in your directory. If the Domain's AuthMode is not SSO, this field cannot be specified.
* `user_settings` - (Required) The user settings. See [User Settings](#user-settings) below.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### User Settings

* `execution_role` - (Required) The execution role ARN for the user.
* `security_groups` - (Optional) The security groups.
* `sharing_settings` - (Optional) The sharing settings. See [Sharing Settings](#sharing-settings) below.
* `tensor_board_app_settings` - (Optional) The TensorBoard app settings. See [TensorBoard App Settings](#tensorboard-app-settings) below.
* `jupyter_server_app_settings` - (Optional) The Jupyter server's app settings. See [Jupyter Server App Settings](#jupyter-server-app-settings) below.
* `kernel_gateway_app_settings` - (Optional) The kernel gateway app settings. See [Kernel Gateway App Settings](#kernel-gateway-app-settings) below.
* `r_session_app_settings` - (Optional) The RSession app settings. See [RSession App Settings](#rsession-app-settings) below.
* `canvas_app_settings` - (Optional) The Canvas app settings. See [Canvas App Settings](#canvas-app-settings) below.

#### Canvas App Settings

* `time_series_forecasting_settings` - (Optional) Time series forecast settings for the Canvas app. see [Time Series Forecasting Settings](#time-series-forecasting-settings) below.

#### Sharing Settings

* `notebook_output_option` - (Optional) Whether to include the notebook cell output when sharing the notebook. The default is `Disabled`. Valid values are `Allowed` and `Disabled`.
* `s3_kms_key_id` - (Optional) When `notebook_output_option` is Allowed, the AWS Key Management Service (KMS) encryption key ID used to encrypt the notebook cell output in the Amazon S3 bucket.
* `s3_output_path` - (Optional) When `notebook_output_option` is Allowed, the Amazon S3 bucket used to save the notebook cell output.

#### TensorBoard App Settings

* `default_resource_spec` - (Optional) The default instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance. see [Default Resource Spec](#default-resource-spec) below.

#### Kernel Gateway App Settings

* `default_resource_spec` - (Optional) The default instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance. see [Default Resource Spec](#default-resource-spec) below.
* `custom_image` - (Optional) A list of custom SageMaker images that are configured to run as a KernelGateway app. see [Custom Image](#custom-image) below.
* `lifecycle_config_arns` - (Optional) The Amazon Resource Name (ARN) of the Lifecycle Configurations.

#### Jupyter Server App Settings

* `code_repository` - (Optional) A list of Git repositories that SageMaker automatically displays to users for cloning in the JupyterServer application. see [Code Repository](#code-repository) below.
* `default_resource_spec` - (Optional) The default instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance. see [Default Resource Spec](#default-resource-spec) below.
* `lifecycle_config_arns` - (Optional) The Amazon Resource Name (ARN) of the Lifecycle Configurations.

#### RSession App Settings

* `default_resource_spec` - (Optional) The default instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance. see [Default Resource Spec](#default-resource-spec) below.
* `custom_image` - (Optional) A list of custom SageMaker images that are configured to run as a KernelGateway app. see [Custom Image](#custom-image) below.

##### Code Repository

* `repository_url` - (Optional) The URL of the Git repository.

##### Default Resource Spec

* `instance_type` - (Optional) The instance type.
* `lifecycle_config_arn` - (Optional) The Amazon Resource Name (ARN) of the Lifecycle Configuration attached to the Resource.
* `sagemaker_image_arn` - (Optional) The Amazon Resource Name (ARN) of the SageMaker image created on the instance.
* `sagemaker_image_version_arn` - (Optional) The ARN of the image version created on the instance.

##### Custom Image

* `app_image_config_name` - (Required) The name of the App Image Config.
* `image_name` - (Required) The name of the Custom Image.
* `image_version_number` - (Optional) The version number of the Custom Image.

##### Time Series Forecasting Settings

* `amazon_forecast_role_arn` - (Optional)  The IAM role that Canvas passes to Amazon Forecast for time series forecasting. By default, Canvas uses the execution role specified in the UserProfile that launches the Canvas app. If an execution role is not specified in the UserProfile, Canvas uses the execution role specified in the Domain that owns the UserProfile. To allow time series forecasting, this IAM role should have the [AmazonSageMakerCanvasForecastAccess](https://docs.aws.amazon.com/sagemaker/latest/dg/security-iam-awsmanpol-canvas.html#security-iam-awsmanpol-AmazonSageMakerCanvasForecastAccess) policy attached and forecast.amazonaws.com added in the trust relationship as a service principal.
* `status` - (Optional) Describes whether time series forecasting is enabled or disabled in the Canvas app. Valid values are `ENABLED` and `DISABLED`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The user profile Amazon Resource Name (ARN).
* `arn` - The user profile Amazon Resource Name (ARN).
* `home_efs_file_system_uid` - The ID of the user's profile in the Amazon Elastic File System (EFS) volume.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

SageMaker User Profiles can be imported using the `arn`, e.g.,

```
$ terraform import aws_sagemaker_user_profile.test_user_profile arn:aws:sagemaker:us-west-2:123456789012:user-profile/domain-id/profile-name
```
