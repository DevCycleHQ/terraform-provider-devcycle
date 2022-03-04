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