---
layout: "selectel"
page_title: "Selectel: selectel_secretsmanager_certificate_v1"
sidebar_current: "docs-selectel-resource-secretsmanager-certificate-v1"
description: |-
    Creates and manages a Certificate in Selectel SecretsManager service using public API v1.
---

# selectel\_secretsmanager\_certificate_v1


Creates and manages a Certificate in Selectel SecretsManager service using public API v1. Learn more about [Certificates](https://docs.selectel.ru/en/cloud/secrets-manager/certificates/).

## Example Usage
```hcl
resource "selectel_secretsmanager_certificate_v1" "cert_1" {
    name = "Terraform-Certificate",
    certificates = [file("./_cert.pem")]
    private_key = file("./_private_key.pem")
    project_id = selectel_vpc_project_v2.project_1.id
}
```

-> It is also possible to pass a EOF sring 
> <details>
> <summary>Expand</summary>
> 
> ```hcl
> resource "selectel_secretsmanager_certificate_v1" "cert_1" {
>     name = "Terraform-Certificate",
>     certificates = [
>         <<-EOF
>         -----BEGIN CERTIFICATE-----
>         MIIDSzCCAjOgAwIBAgIULEumDHpDEHvQ1seZB9yRX9sCgoUwDQYJKoZIhvcNAQEL
>         ...
>         ----END CERTIFICATE-----
>         EOF
>     ]
>     private_key = <<-EOF
>     -----BEGIN PRIVATE KEY-----
>     MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCuk3SFn0AfAoxo
>     ...
>     -----END PRIVATE KEY-----
>     EOF
>     project_id = selectel_vpc_project_v2.project_1.id
> }
> ```
> </details>

## Argument Reference
- `name` (Required) — name of the certificate.
- `certificates` (Required) — ca_chain list of public certificates for certificate. Each certificate must begin with **-----BEGIN CERTIFICATE-----** and end with **-----END CERTIFICATE-----**.
- `private_key` (Required, Sensitive) — that should start with **-----BEGIN PRIVATE KEY-----** and end with **-----END PRIVATE KEY-----**.
- `project_id` (Required) — unique identifier of the associated Cloud Platform project.

## Attributes Reference
- `dns_names` — domain names.
- `id` — computed id of a certificate.
- `issued_by` — information that is incorporated into certificate.
- `serial` — number written in the certificate that was chosen by the CA which issued the certificate.
- `validity` — validity of a certificate in terms of notBefore and notAfter timestamps.
- `version` — of the certificate.

## Import

~> When importing Certificate you have to provide unique identifier of the associated Cloud Platform project

### Using import block
-> In Terraform v1.5.0 and later, use an import block to import Certificate using template below.

```hcl
import {
   to = selectel_secretsmanager_certificate_v1.imported_certificate
   id = "<id>"
}
```

* `<id>` — Unique identifier of the certificate. To get the id of the certificate in the [Control panel](https://my.selectel.ru/vpc/), go to **Cloud Platform** ⟶ **Secrets Manager** ⟶ **Certificate** in contect menu click copy UUID.


### Using terraform import
```shell
export SEL_PROJECT_ID=<selectel_project_id>
terraform import selectel_secretsmanager_certificate_v1.imported_certificate <id>
```

* `<selectel_project_id>` — Unique identifier of the associated Cloud Platform project. To get the project ID, in the [Control panel](https://my.selectel.ru/vpc/), go to **Cloud Platform** ⟶ project name ⟶ copy the ID of the required project. Learn more about [Cloud Platform projects](https://docs.selectel.ru/cloud/managed-databases/about/projects/).

* `<id>` — Unique identifier of the certificate. To get the id of the certificate in the [Control panel](https://my.selectel.ru/vpc/), go to **Cloud Platform** ⟶ **Secrets Manager** ⟶ **Certificate** in contect menu click copy UUID.


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
│   with selectel_secretsmanager_certificate_v1.imported_certificate,
│   on generated_resources.tf line 1:
│   (source code not available)
│ 
│ The argument "certificates" is required, but no definition was found.
╵
╷
│ Error: Missing required argument
│ 
│   with selectel_secretsmanager_certificate_v1.imported_certificate,
│   on generated_resources.tf line 3:
│   (source code not available)
│ 
│ The argument "private_key" is required, but no definition was found.
```

~> However `generated_resources.tf` was correctly generated:

```hcl
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform
resource "selectel_secretsmanager_certificate_v1" "imported_certificate" {
  certificates = null
  name         = "Cert-from-Cloud"
  private_key  = null # sensitive
  project_id   = <selectel_project_id>
}
```

All you have to do is to set `certificates` attribute from `null` to `[]` and `private_key` attribute from `null` to `"null"`. This move doesn't Destroy and Create certificate.
```hcl
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.
 
# __generated__ by Terraform
resource "selectel_secretsmanager_certificate_v1" "imported_certificate" {
  certificates = []                     # <- set null to []
  name         = "Cert-from-Cloud"
  private_key  = "null" # sensitive     # <- set null to string
  project_id   = <selectel_project_id>
}
```