---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "devcycle_feature Data Source - terraform-provider-devcycle"
subcategory: ""
description: |-
  DevCycle Feature data source
---

# devcycle_feature (Data Source)

DevCycle Feature data source

## Example Usage

```terraform
data "devcycle_feature" "test" {
  key         = "terraform-provider-feature"
  project_key = "terraform-provider-testing"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **key** (String) Feature key
- **project_key** (String) Project key that the feature belongs to

### Optional

- **variables** (Attributes List) Feature variables (see [below for nested schema](#nestedatt--variables))
- **variations** (Attributes List) Feature variations (see [below for nested schema](#nestedatt--variations))

### Read-Only

- **description** (String) Feature description
- **id** (String) Feature ID
- **name** (String) Feature name
- **project_id** (String) Project ID that the feature belongs to
- **type** (String) Feature Type

<a id="nestedatt--variables"></a>
### Nested Schema for `variables`

Optional:

- **created_at** (String) Created at timestamp
- **description** (String) Variation feature key
- **feature_key** (String) Variation feature key
- **id** (String) Variation type
- **key** (String) Variation key
- **name** (String) Variation name
- **type** (String) Variation type
- **updated_at** (String) Updated at timestamp


<a id="nestedatt--variations"></a>
### Nested Schema for `variations`

Optional:

- **id** (String) Variation type
- **key** (String) Variation key
- **name** (String) Variation name
- **variables** (Map of String) Variation variables


