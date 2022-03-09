data "devcycle_evaluated_variable_boolean" "test" {
  id = "acceptance-testing-boolean"
  user = {
    id = "acceptancetesting"
  }
  default_value = false
}