data "devcycle_evaluated_variable_number" "test" {
  id = "acceptance-testing-number"
  user = {
    id = "acceptancetesting"
  }
  default_value = 1
}