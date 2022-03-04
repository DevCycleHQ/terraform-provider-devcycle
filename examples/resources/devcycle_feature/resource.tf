resource "devcycle_feature" "test" {
  project_id  = "622112634cabe0e9fbaf974d"
  name        = "TerraformAccTest"
  key         = "terraform-acceptance-testing"
  description = "Terraform acceptance testing"
  type        = "experiment"
  tags        = ["acceptance-testing"]
}