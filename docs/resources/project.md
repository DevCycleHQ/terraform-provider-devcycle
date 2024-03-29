---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "devcycle_project Resource - terraform-provider-devcycle"
subcategory: ""
description: |-
  DevCycle project resource. Allows for creation/modification of a project.
---

# devcycle_project (Resource)

DevCycle project resource. Allows for creation/modification of a project.

## Example Usage

```terraform
resource "devcycle_project" "test" {
  name        = "TerraformAccTest"
  key         = "project-key"
  description = "Terraform acceptance testing"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Description of the project
- `key` (String) Project key, usually the lowercase, kebab case name of the project
- `name` (String) Name of the project

### Read-Only

- `id` (String) Project Id
- `organization` (String) Organization that the project belongs to


