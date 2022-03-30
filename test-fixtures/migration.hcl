
migration "state" "20220331032518_tfmigration" {
  dir       = ""
  actions   = ["mv \"aws_security_group.foo\" \"aws_security_group.baz\""]
  force     = false
  workspace = ""
}
