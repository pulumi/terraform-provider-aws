---
subcategory: "CodeBuild"
layout: "aws"
page_title: "AWS: aws_codebuild_project"
description: |-
  Provides a CodeBuild Project resource.
---

# Resource: aws_codebuild_project

Provides a CodeBuild Project resource. See also the `aws_codebuild_webhook` resource, which manages the webhook to the source (e.g., the "rebuild every time a code change is pushed" option in the CodeBuild web console).

## Example Usage

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "example"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

resource "aws_iam_role" "example" {
  name = "example"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "example" {
  role = aws_iam_role.example.name

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:CreateNetworkInterface",
        "ec2:DescribeDhcpOptions",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DeleteNetworkInterface",
        "ec2:DescribeSubnets",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeVpcs"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:CreateNetworkInterfacePermission"
      ],
      "Resource": [
        "arn:aws:ec2:us-east-1:123456789012:network-interface/*"
      ],
      "Condition": {
        "StringEquals": {
          "ec2:Subnet": [
            "${aws_subnet.example1.arn}",
            "${aws_subnet.example2.arn}"
          ],
          "ec2:AuthorizedService": "codebuild.amazonaws.com"
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "${aws_s3_bucket.example.arn}",
        "${aws_s3_bucket.example.arn}/*"
      ]
    }
  ]
}
POLICY
}

resource "aws_codebuild_project" "example" {
  name          = "test-project"
  description   = "test_codebuild_project"
  build_timeout = "5"
  service_role  = aws_iam_role.example.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  cache {
    type     = "S3"
    location = aws_s3_bucket.example.bucket
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:1.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "SOME_KEY1"
      value = "SOME_VALUE1"
    }

    environment_variable {
      name  = "SOME_KEY2"
      value = "SOME_VALUE2"
      type  = "PARAMETER_STORE"
    }
  }

  logs_config {
    cloudwatch_logs {
      group_name  = "log-group"
      stream_name = "log-stream"
    }

    s3_logs {
      status   = "ENABLED"
      location = "${aws_s3_bucket.example.id}/build-log"
    }
  }

  source {
    type            = "GITHUB"
    location        = "https://github.com/mitchellh/packer.git"
    git_clone_depth = 1

    git_submodules_config {
      fetch_submodules = true
    }
  }

  source_version = "master"

  vpc_config {
    vpc_id = aws_vpc.example.id

    subnets = [
      aws_subnet.example1.id,
      aws_subnet.example2.id,
    ]

    security_group_ids = [
      aws_security_group.example1.id,
      aws_security_group.example2.id,
    ]
  }

  tags = {
    Environment = "Test"
  }
}

resource "aws_codebuild_project" "project-with-cache" {
  name           = "test-project-cache"
  description    = "test_codebuild_project_cache"
  build_timeout  = "5"
  queued_timeout = "5"

  service_role = aws_iam_role.example.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  cache {
    type  = "LOCAL"
    modes = ["LOCAL_DOCKER_LAYER_CACHE", "LOCAL_SOURCE_CACHE"]
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:1.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "SOME_KEY1"
      value = "SOME_VALUE1"
    }
  }

  source {
    type            = "GITHUB"
    location        = "https://github.com/mitchellh/packer.git"
    git_clone_depth = 1
  }

  tags = {
    Environment = "Test"
  }
}
```

## Argument Reference

The following arguments are required:

* `artifacts` - (Required) Configuration block. Detailed below.
* `environment` - (Required) Configuration block. Detailed below.
* `name` - (Required) Project's name.
* `service_role` - (Required) Amazon Resource Name (ARN) of the AWS Identity and Access Management (IAM) role that enables AWS CodeBuild to interact with dependent AWS services on behalf of the AWS account.
* `source` - (Required) Configuration block. Detailed below.

The following arguments are optional:

* `badge_enabled` - (Optional) Generates a publicly-accessible URL for the projects build badge. Available as `badge_url` attribute when enabled.
* `build_batch_config` - (Optional) Defines the batch build options for the project.
* `build_timeout` - (Optional) Number of minutes, from 5 to 480 (8 hours), for AWS CodeBuild to wait until timing out any related build that does not get marked as completed. The default is 60 minutes.
* `cache` - (Optional) Configuration block. Detailed below.
* `concurrent_build_limit` - (Optional) Specify a maximum number of concurrent builds for the project. The value specified must be greater than 0 and less than the account concurrent running builds limit.
* `description` - (Optional) Short description of the project.
* `file_system_locations` - (Optional) A set of file system locations to mount inside the build. File system locations are documented below.
* `encryption_key` - (Optional) AWS Key Management Service (AWS KMS) customer master key (CMK) to be used for encrypting the build project's build output artifacts.
* `logs_config` - (Optional) Configuration block. Detailed below.
* `project_visibility` - (Optional) Specifies the visibility of the project's builds. Possible values are: `PUBLIC_READ` and `PRIVATE`. Default value is `PRIVATE`.
* `resource_access_role` - The ARN of the IAM role that enables CodeBuild to access the CloudWatch Logs and Amazon S3 artifacts for the project's builds.
* `queued_timeout` - (Optional) Number of minutes, from 5 to 480 (8 hours), a build is allowed to be queued before it times out. The default is 8 hours.
* `secondary_artifacts` - (Optional) Configuration block. Detailed below.
* `secondary_sources` - (Optional) Configuration block. Detailed below.
* `secondary_source_version` - (Optional) Configuration block. Detailed below.
* `source_version` - (Optional) Version of the build input to be built for this project. If not specified, the latest version is used.
* `tags` - (Optional) Map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `vpc_config` - (Optional) Configuration block. Detailed below.

### artifacts

* `artifact_identifier` - (Optional) Artifact identifier. Must be the same specified inside the AWS CodeBuild build specification.
* `bucket_owner_access` - (Optional) Specifies the bucket owner's access for objects that another account uploads to their Amazon S3 bucket. By default, only the account that uploads the objects to the bucket has access to these objects. This property allows you to give the bucket owner access to these objects. Valid values are `NONE`, `READ_ONLY`, and `FULL`. your CodeBuild service role must have the `s3:PutBucketAcl` permission. This permission allows CodeBuild to modify the access control list for the bucket.
* `encryption_disabled` - (Optional) Whether to disable encrypting output artifacts. If `type` is set to `NO_ARTIFACTS`, this value is ignored. Defaults to `false`.
* `location` - (Optional) Information about the build output artifact location. If `type` is set to `CODEPIPELINE` or `NO_ARTIFACTS`, this value is ignored. If `type` is set to `S3`, this is the name of the output bucket.
* `name` - (Optional) Name of the project. If `type` is set to `S3`, this is the name of the output artifact object
* `namespace_type` - (Optional) Namespace to use in storing build artifacts. If `type` is set to `S3`, then valid values are `BUILD_ID`, `NONE`.
* `override_artifact_name` (Optional) Whether a name specified in the build specification overrides the artifact name.
* `packaging` - (Optional) Type of build output artifact to create. If `type` is set to `S3`, valid values are `NONE`, `ZIP`
* `path` - (Optional) If `type` is set to `S3`, this is the path to the output artifact.
* `type` - (Required) Build output artifact's type. Valid values: `CODEPIPELINE`, `NO_ARTIFACTS`, `S3`.

### build_batch_config

* `combine_artifacts` - (Optional) Specifies if the build artifacts for the batch build should be combined into a single artifact location.
* `restrictions` - (Optional) Configuration block specifying the restrictions for the batch build. Detailed below.
* `service_role` - (Required) Specifies the service role ARN for the batch build project.
* `timeout_in_mins` - (Optional) Specifies the maximum amount of time, in minutes, that the batch build must be completed in.

#### build_batch_config: restrictions

* `compute_types_allowed` - (Optional) An array of strings that specify the compute types that are allowed for the batch build. See [Build environment compute types](https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-compute-types.html) in the AWS CodeBuild User Guide for these values.
* `maximum_builds_allowed` - (Optional) Specifies the maximum number of builds allowed.

### cache

* `location` - (Required when cache type is `S3`) Location where the AWS CodeBuild project stores cached resources. For type `S3`, the value must be a valid S3 bucket name/prefix.
* `modes` - (Required when cache type is `LOCAL`) Specifies settings that AWS CodeBuild uses to store and reuse build dependencies. Valid values:  `LOCAL_SOURCE_CACHE`, `LOCAL_DOCKER_LAYER_CACHE`, `LOCAL_CUSTOM_CACHE`.
* `type` - (Optional) Type of storage that will be used for the AWS CodeBuild project cache. Valid values: `NO_CACHE`, `LOCAL`, `S3`. Defaults to `NO_CACHE`.

### environment

* `certificate` - (Optional) ARN of the S3 bucket, path prefix and object key that contains the PEM-encoded certificate.
* `compute_type` - (Required) Information about the compute resources the build project will use. Valid values: `BUILD_GENERAL1_SMALL`, `BUILD_GENERAL1_MEDIUM`, `BUILD_GENERAL1_LARGE`, `BUILD_GENERAL1_2XLARGE`. `BUILD_GENERAL1_SMALL` is only valid if `type` is set to `LINUX_CONTAINER`. When `type` is set to `LINUX_GPU_CONTAINER`, `compute_type` must be `BUILD_GENERAL1_LARGE`.
* `environment_variable` - (Optional) Configuration block. Detailed below.
* `image_pull_credentials_type` - (Optional) Type of credentials AWS CodeBuild uses to pull images in your build. Valid values: `CODEBUILD`, `SERVICE_ROLE`. When you use a cross-account or private registry image, you must use SERVICE_ROLE credentials. When you use an AWS CodeBuild curated image, you must use CodeBuild credentials. Defaults to `CODEBUILD`.
* `image` - (Required) Docker image to use for this build project. Valid values include [Docker images provided by CodeBuild](https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-available.html) (e.g `aws/codebuild/standard:2.0`), [Docker Hub images](https://hub.docker.com/) (e.g., `nginx/nginx:latest`), and full Docker repository URIs such as those for ECR (e.g., `137112412989.dkr.ecr.us-west-2.amazonaws.com/amazonlinux:latest`).
* `privileged_mode` - (Optional) Whether to enable running the Docker daemon inside a Docker container. Defaults to `false`.
* `registry_credential` - (Optional) Configuration block. Detailed below.
* `type` - (Required) Type of build environment to use for related builds. Valid values: `LINUX_CONTAINER`, `LINUX_GPU_CONTAINER`, `WINDOWS_CONTAINER` (deprecated), `WINDOWS_SERVER_2019_CONTAINER`, `ARM_CONTAINER`. For additional information, see the [CodeBuild User Guide](https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-compute-types.html).

#### environment: environment_variable

* `name` - (Required) Environment variable's name or key.
* `type` - (Optional) Type of environment variable. Valid values: `PARAMETER_STORE`, `PLAINTEXT`, `SECRETS_MANAGER`.
* `value` - (Required) Environment variable's value.

#### environment: registry_credential

Credentials for access to a private Docker registry.

* `credential` - (Required) ARN or name of credentials created using AWS Secrets Manager.
* `credential_provider` - (Required) Service that created the credentials to access a private Docker registry. Valid value: `SECRETS_MANAGER` (AWS Secrets Manager).

### file_system_locations

See [ProjectFileSystemLocation](https://docs.aws.amazon.com/codebuild/latest/APIReference/API_ProjectFileSystemLocation.html) for more details of the fields.

* `identifier` - (Optional) The name used to access a file system created by Amazon EFS. CodeBuild creates an environment variable by appending the identifier in all capital letters to CODEBUILD\_. For example, if you specify my-efs for identifier, a new environment variable is create named CODEBUILD_MY-EFS.
* `location` - (Optional) A string that specifies the location of the file system created by Amazon EFS. Its format is `efs-dns-name:/directory-path`.
* `mount_options` - (Optional) The mount options for a file system created by AWS EFS.
* `mount_point` - (Optional) The location in the container where you mount the file system.
* `type` - (Optional) The type of the file system. The one supported type is `EFS`.

### logs_config

* `cloudwatch_logs` - (Optional) Configuration block. Detailed below.
* `s3_logs` - (Optional) Configuration block. Detailed below.

#### logs_config: cloudwatch_logs

* `group_name` - (Optional) Group name of the logs in CloudWatch Logs.
* `status` - (Optional) Current status of logs in CloudWatch Logs for a build project. Valid values: `ENABLED`, `DISABLED`. Defaults to `ENABLED`.
* `stream_name` - (Optional) Stream name of the logs in CloudWatch Logs.

#### logs_config: s3_logs

* `encryption_disabled` - (Optional) Whether to disable encrypting S3 logs. Defaults to `false`.
* `location` - (Optional) Name of the S3 bucket and the path prefix for S3 logs. Must be set if status is `ENABLED`, otherwise it must be empty.
* `status` - (Optional) Current status of logs in S3 for a build project. Valid values: `ENABLED`, `DISABLED`. Defaults to `DISABLED`.
* `bucket_owner_access` - (Optional) Specifies the bucket owner's access for objects that another account uploads to their Amazon S3 bucket. By default, only the account that uploads the objects to the bucket has access to these objects. This property allows you to give the bucket owner access to these objects. Valid values are `NONE`, `READ_ONLY`, and `FULL`. your CodeBuild service role must have the `s3:PutBucketAcl` permission. This permission allows CodeBuild to modify the access control list for the bucket.

### secondary_artifacts

* `artifact_identifier` - (Required) Artifact identifier. Must be the same specified inside the AWS CodeBuild build specification.
* `bucket_owner_access` - (Optional) Specifies the bucket owner's access for objects that another account uploads to their Amazon S3 bucket. By default, only the account that uploads the objects to the bucket has access to these objects. This property allows you to give the bucket owner access to these objects. Valid values are `NONE`, `READ_ONLY`, and `FULL`. your CodeBuild service role must have the `s3:PutBucketAcl` permission. This permission allows CodeBuild to modify the access control list for the bucket.
* `encryption_disabled` - (Optional) Whether to disable encrypting output artifacts. If `type` is set to `NO_ARTIFACTS`, this value is ignored. Defaults to `false`.
* `location` - (Optional) Information about the build output artifact location. If `type` is set to `CODEPIPELINE` or `NO_ARTIFACTS`, this value is ignored. If `type` is set to `S3`, this is the name of the output bucket. If `path` is not also specified, then `location` can also specify the path of the output artifact in the output bucket.
* `name` - (Optional) Name of the project. If `type` is set to `S3`, this is the name of the output artifact object
* `namespace_type` - (Optional) Namespace to use in storing build artifacts. If `type` is set to `S3`, then valid values are `BUILD_ID` or `NONE`.
* `override_artifact_name` (Optional) Whether a name specified in the build specification overrides the artifact name.
* `packaging` - (Optional) Type of build output artifact to create. If `type` is set to `S3`, valid values are `NONE`, `ZIP`
* `path` - (Optional) If `type` is set to `S3`, this is the path to the output artifact.
* `type` - (Required) Build output artifact's type. The only valid value is `S3`.

### secondary_sources

* `auth` - (Optional, **Deprecated**) Configuration block with the authorization settings for AWS CodeBuild to access the source code to be built. This information is for the AWS CodeBuild console's use only. Use the `aws_codebuild_source_credential` resource instead. Auth blocks are documented below.
* `buildspec` - (Optional) The build spec declaration to use for this build project's related builds. This must be set when `type` is `NO_SOURCE`. It can either be a path to a file residing in the repository to be built or a local file path leveraging the `file()` built-in.
* `git_clone_depth` - (Optional) Truncate git history to this many commits. Use `0` for a `Full` checkout which you need to run commands like `git branch --show-current`. See [AWS CodePipeline User Guide: Tutorial: Use full clone with a GitHub pipeline source](https://docs.aws.amazon.com/codepipeline/latest/userguide/tutorials-github-gitclone.html) for details.
* `git_submodules_config` - (Optional) Configuration block. Detailed below.
* `insecure_ssl` - (Optional) Ignore SSL warnings when connecting to source control.
* `location` - (Optional) Location of the source code from git or s3.
* `report_build_status` - (Optional) Whether to report the status of a build's start and finish to your source provider. This option is only valid when your source provider is `GITHUB`, `BITBUCKET`, or `GITHUB_ENTERPRISE`.
* `build_status_config` - (Optional) Configuration block that contains information that defines how the build project reports the build status to the source provider. This option is only used when the source provider is `GITHUB`, `GITHUB_ENTERPRISE`, or `BITBUCKET`. `build_status_config` blocks are documented below.
* `source_identifier` - (Required) An identifier for this project source. The identifier can only contain alphanumeric characters and underscores, and must be less than 128 characters in length.
* `type` - (Required) Type of repository that contains the source code to be built. Valid values: `CODECOMMIT`, `CODEPIPELINE`, `GITHUB`, `GITHUB_ENTERPRISE`, `BITBUCKET` or `S3`.

#### secondary_sources: auth

* `resource` - (Optional, **Deprecated**) Resource value that applies to the specified authorization type. Use the `aws_codebuild_source_credential` resource instead.
* `type` - (Required, **Deprecated**) Authorization type to use. The only valid value is `OAUTH`. This data type is deprecated and is no longer accurate or used. Use the `aws_codebuild_source_credential` resource instead.

#### secondary_sources: git_submodules_config

This block is only valid when the `type` is `CODECOMMIT`, `GITHUB` or `GITHUB_ENTERPRISE`.

* `fetch_submodules` - (Required) Whether to fetch Git submodules for the AWS CodeBuild build project.

#### secondary_sources: build_status_config

* `context` - (Optional) Specifies the context of the build status CodeBuild sends to the source provider. The usage of this parameter depends on the source provider.
* `target_url` - (Optional) Specifies the target url of the build status CodeBuild sends to the source provider. The usage of this parameter depends on the source provider.

### secondary_source_version

* `source_identifier` - (Required) An identifier for a source in the build project.
* `source_version` - (Required) The source version for the corresponding source identifier. See [AWS docs](https://docs.aws.amazon.com/codebuild/latest/APIReference/API_ProjectSourceVersion.html#CodeBuild-Type-ProjectSourceVersion-sourceVersion) for more details.

### source

* `auth` - (Optional, **Deprecated**) Configuration block with the authorization settings for AWS CodeBuild to access the source code to be built. This information is for the AWS CodeBuild console's use only. Use the `aws_codebuild_source_credential` resource instead. Auth blocks are documented below.
* `buildspec` - (Optional) Build specification to use for this build project's related builds. This must be set when `type` is `NO_SOURCE`.
* `git_clone_depth` - (Optional) Truncate git history to this many commits. Use `0` for a `Full` checkout which you need to run commands like `git branch --show-current`. See [AWS CodePipeline User Guide: Tutorial: Use full clone with a GitHub pipeline source](https://docs.aws.amazon.com/codepipeline/latest/userguide/tutorials-github-gitclone.html) for details.
* `git_submodules_config` - (Optional) Configuration block. Detailed below.
* `insecure_ssl` - (Optional) Ignore SSL warnings when connecting to source control.
* `location` - (Optional) Location of the source code from git or s3.
* `report_build_status` - (Optional) Whether to report the status of a build's start and finish to your source provider. This option is only valid when the `type` is `BITBUCKET` or `GITHUB`.
* `build_status_config` - (Optional) Configuration block that contains information that defines how the build project reports the build status to the source provider. This option is only used when the source provider is `GITHUB`, `GITHUB_ENTERPRISE`, or `BITBUCKET`. `build_status_config` blocks are documented below.
* `type` - (Required) Type of repository that contains the source code to be built. Valid values: `CODECOMMIT`, `CODEPIPELINE`, `GITHUB`, `GITHUB_ENTERPRISE`, `BITBUCKET`, `S3`, `NO_SOURCE`.

#### source: auth

* `resource` - (Optional, **Deprecated**) Resource value that applies to the specified authorization type. Use the `aws_codebuild_source_credential` resource instead.
* `type` - (Required, **Deprecated**) Authorization type to use. The only valid value is `OAUTH`. This data type is deprecated and is no longer accurate or used. Use the `aws_codebuild_source_credential` resource instead.

#### source: git_submodules_config

This block is only valid when the `type` is `CODECOMMIT`, `GITHUB` or `GITHUB_ENTERPRISE`.

* `fetch_submodules` - (Required) Whether to fetch Git submodules for the AWS CodeBuild build project.

#### source: build_status_config

* `context` - (Optional) Specifies the context of the build status CodeBuild sends to the source provider. The usage of this parameter depends on the source provider.
* `target_url` - (Optional) Specifies the target url of the build status CodeBuild sends to the source provider. The usage of this parameter depends on the source provider.

### vpc_config

* `security_group_ids` - (Required) Security group IDs to assign to running builds.
* `subnets` - (Required) Subnet IDs within which to run builds.
* `vpc_id` - (Required) ID of the VPC within which to run builds.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the CodeBuild project.
* `badge_url` - URL of the build badge when `badge_enabled` is enabled.
* `id` - Name (if imported via `name`) or ARN (if created via the provider or imported via ARN) of the CodeBuild project.
* `public_project_alias` - The project identifier used with the public build APIs.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

CodeBuild Project can be imported using the `name`, e.g.,

```
$ terraform import aws_codebuild_project.name project-name
```
