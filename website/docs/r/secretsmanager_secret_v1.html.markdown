---
layout: "selectel"
page_title: "Selectel: selectel_secretsmanager_secret_v1"
sidebar_current: "docs-selectel-resource-secretsmanager-secret-v1"
description: |-
    Creates and manages a Secret in Selectel SecretsManager service using public API v1.
---

# selectel\_secretsmanager\_secret_v1

Creates and manages a Secret in Selectel SecretsManager service using public API v1. Learn more about [Secrets](https://docs.selectel.ru/en/cloud/secrets-manager/secrets/).

## Example Usage
```hcl
resource "selectel_secretsmanager_secret_v1" "secret_1" {
    key = "Terraform-Secret"
    description = "Secret from .tf"
    value = "zelibobs"
    project_id = selectel_vpc_project_v2.tf_secretsmanager.id
}
```

## Argument Reference
- `key` (Required) — unique key, name of the secret.
- `description` (Optinal) — description of the secret.
- `value` (Required, Sensitive) — secret value, e.g. password, API key, certificate key, or other.
- `project_id` (Required) — unique identifier of the associated Cloud Platform project.

## Attributes Reference
- `created_at` — time when the secret was created.
- `name` — computed name of the secret same as key.

## Import

~> When importing Secret you have to provide unique identifier of the associated Cloud Platform project

### Using import block
-> In Terraform v1.5.0 and later, use an import block to import Secret using template below.

```hcl
import {
   to = selectel_secretsmanager_secret_v1.imported_secret
   id = <key>
}
```

* `<key>` — Unique identifier of the secret, its name. To get the name of the secret in the [Control panel](https://my.selectel.ru/vpc/), go to **Cloud Platform** ⟶ **Secrets Manager** ⟶ **Secret** copy the Name.



### Using terraform import
```shell
export SEL_PROJECT_ID=<selectel_project_id>
terraform import selectel_secretsmanager_secret_v1.imported_secret <key>
```

* `<selectel_project_id>` — Unique identifier of the associated Cloud Platform project. To get the project ID, in the [Control panel](https://my.selectel.ru/vpc/), go to **Cloud Platform** ⟶ project name ⟶ copy the ID of the required project. Learn more about [Cloud Platform projects](https://docs.selectel.ru/cloud/managed-databases/about/projects/).

* `<key>` — Unique identifier of the secret, its name. To get the name of the secret in the [Control panel](https://my.selectel.ru/vpc/), go to **Cloud Platform** ⟶ **Secrets Manager** ⟶ **Secret** copy the Name.


!> Generating configuration
Because, configuration generation is available in Terraform v1.5 as an experimental feature. It struggles while setting Sensetive and Required

```shell
SEL_PROJECT_ID=<selectel_project_id> terraform plan -generate-config-out=generated_resources.tf
```

```text
Planning failed. Terraform encountered an error while generating this plan.
 
╷
│ Error: Missing required argument
│
│   with selectel_secretsmanager_secret_v1.import_secret,
│   on generated_resources.tf line 4:
│   (source code not available)
│
│ The argument "value" is required, but no definition was found.
```

~> However `generated_resources.tf` was correctly generated:

```hcl
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.
 
# __generated__ by Terraform
resource "selectel_secretsmanager_secret_v1" "import_secret" {
  description = "123456"
  key         = "Secret-from-Cloud"
  project_id  = <selectel_project_id>
  value       = null # sensitive
}
```

All you have to do is to set `null` to `"null"` for example, this move doesn't Destroy and Create secret.
```hcl
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.
 
# __generated__ by Terraform
resource "selectel_secretsmanager_secret_v1" "import_secret" {
  description = "123456"
  key         = "Secret-from-Cloud"
  project_id  = <selectel_project_id>
  value       = "null" # sensitive    # <- set null to string
}
```