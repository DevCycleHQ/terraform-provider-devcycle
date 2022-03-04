resource "devcycle_project" "demo" {
  name        = "DemoProject"
  key         = "terraform-demo-project"
  description = "Terraform provider demo"
}

resource "devcycle_environment" "demo" {
  project_id  = devcycle_project.demo.id
  name        = "Terraform Demo Environment"
  key         = "terraform-demo-environment"
  description = "Terraform provider demo"
  color       = "#232323"
  type        = "development"
  settings = {
    app_icon_uri = "demo"
  }
}

resource "devcycle_feature" "demo" {
  project_id  = devcycle_project.demo.id
  name        = "Terraform Demo Feature"
  key         = "terraform-demo-feature"
  description = "Terraform provider demo"
  type        = "experiment"
  tags        = ["terraform-demo"]
}

resource "devcycle_variable" "test" {
  name          = "Terraform Demo Variable"
  key           = "terraform-acc-test"
  description   = "Terraform acceptance testing"
  type          = "Boolean"
  feature_id    = devcycle_feature.demo.id
  project_id    = devcycle_project.demo.id
  default_value = "false"
}