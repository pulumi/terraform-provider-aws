---
subcategory: "Cognito Identity"
layout: "aws"
page_title: "AWS: aws_cognito_identity_pool_roles_attachment"
description: |-
  Provides an AWS Cognito Identity Pool Roles Attachment.
---

# Resource: aws_cognito_identity_pool_roles_attachment

Provides an AWS Cognito Identity Pool Roles Attachment.

## Example Usage


## Argument Reference

The Cognito Identity Pool Roles Attachment argument layout is a structure composed of several sub-resources - these resources are laid out below.

* `identity_pool_id` (Required) - An identity pool ID in the format `REGION_GUID`.
* `role_mapping` (Optional) - A List of [Role Mapping](#role-mappings).
* `roles` (Required) - The map of roles associated with this pool. For a given role, the key will be either "authenticated" or "unauthenticated" and the value will be the Role ARN.

#### Role Mappings

* `identity_provider` (Required) - A string identifying the identity provider, for example, "graph.facebook.com" or "cognito-idp.us-east-1.amazonaws.com/us-east-1_abcdefghi:app_client_id". Depends on `cognito_identity_providers` set on `aws_cognito_identity_pool` resource or a `aws_cognito_identity_provider` resource.
* `ambiguous_role_resolution` (Optional) - Specifies the action to be taken if either no rules match the claim value for the Rules type, or there is no cognito:preferred_role claim and there are multiple cognito:roles matches for the Token type. `Required` if you specify Token or Rules as the Type.
* `mapping_rule` (Optional) - The [Rules Configuration](#rules-configuration) to be used for mapping users to roles. You can specify up to 25 rules per identity provider. Rules are evaluated in order. The first one to match specifies the role.
* `type` (Required) - The role mapping type.

#### Rules Configuration

* `claim` (Required) - The claim name that must be present in the token, for example, "isAdmin" or "paid".
* `match_type` (Required) - The match condition that specifies how closely the claim value in the IdP token must match Value.
* `role_arn` (Required) - The role ARN.
* `value` (Required) - A brief string that the claim must match, for example, "paid" or "yes".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identity pool ID.

## Import

Cognito Identity Pool Roles Attachment can be imported using the Identity Pool ID, e.g.,

```
$ terraform import aws_cognito_identity_pool_roles_attachment.example us-west-2:b64805ad-cb56-40ba-9ffc-f5d8207e6d42
```
