resource "devcycle_variable" "test" {
  name          = "TerraformAccTest"
  key           = "terraform-acc-test"
  description   = "Terraform acceptance testing"
  type          = "Boolean"
  feature_id    = "622115014b06357d06d1cf3e"
  project_id    = "622112634cabe0e9fbaf974d"
  default_value = "false"
}