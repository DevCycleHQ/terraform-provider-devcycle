data "devcycle_evaluated_variable_json" "test" {
  id = "acceptance-testing-json"
  user = {
    id = "acceptancetesting"
  }
  default_value = "{}"
}