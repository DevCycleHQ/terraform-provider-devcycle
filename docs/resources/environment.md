---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "devcycle_environment Resource - terraform-provider-devcycle"
subcategory: ""
description: |-
  DevCycle Environment resource. This resource is used to create and manage DevCycle environments.
---

# devcycle_environment (Resource)

DevCycle Environment resource. This resource is used to create and manage DevCycle environments.

## Example Usage

```terraform
resource "devcycle_environment" "test" {
  project_id  = "project_id"
  name        = "TerraformAccTest"
  key         = "terraform-acceptance-testing"
  description = "Terraform acceptance testing"
  color       = "#232323"
  type        = "development"
  settings = {
    app_icon_uri = "test"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `color` (String) Environment Color in Hex with leading #
- `description` (String) Environment Description
- `key` (String) Environment Key
- `name` (String) Environment Name
- `project_id` (String) Project id or key of the project to which the environment belongs. Using the key (human readable name) is recommended when not managing the project through Terraform.
- `settings` (Attributes) Environment Settings (see [below for nested schema](#nestedatt--settings))
- `type` (String) Environment Type

### Read-Only

- `id` (String) Environment Id
- `sdk_keys` (List of String) SDK Keys for the environment

<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Required:

- `app_icon_uri` (String) Environment App Icon Uri


