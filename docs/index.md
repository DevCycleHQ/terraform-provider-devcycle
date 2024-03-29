---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "devcycle Provider"
subcategory: ""
description: |-
  This provider allows you to manage DevCycle projects, environments, features, and variables. It uses the DevCycle API to manage these resources.  You can find more information about the DevCycle API here https://docs.devcycle.com/management-api/.
  This provider is compatible with Terraform v1.0 and newer. Because of the way that authentication for the management api works - this provider will have access to manage all projects within a DevCycle org. Be careful!
---

# devcycle Provider

This provider allows you to manage DevCycle projects, environments, features, and variables. It uses the DevCycle API to manage these resources.  You can find more information about the DevCycle API [here](https://docs.devcycle.com/management-api/).

This provider is compatible with Terraform v1.0 and newer. Because of the way that authentication for the management api works - this provider will have access to manage all projects within a DevCycle org. Be careful!

## Example Usage

```terraform
provider "devcycle" {
  client_id     = "Client ID from DevCycle, or DEVCYCLE_CLIENT_ID environment variable"
  client_secret = "Client Secret from DevCycle, or DEVCYCLE_CLIENT_SECRET environment variable"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `client_id` (String, Sensitive) API Authentication Client ID. Found in your DevCycle account settings.
- `client_secret` (String, Sensitive) API Authentication Client Secret. Found in your DevCycle account settings.
- `server_sdk_token` (String, Sensitive) Server SDK Token. This is specific to a given project, and an environment. Used to identify and authenticate server sdk requests to evaluate feature flags.
