data "devcycle_evaluated_variable_string" "test" {
  id = "acceptance-testing-string"
  user = {
    id = "acceptancetesting"
  }
  default_value = "string"
}