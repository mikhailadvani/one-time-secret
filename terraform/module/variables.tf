variable region {
  description = "AWS region to create the components in."
}

variable stack_name {
  default = "one-time-secret"
  description = "Name to the lambda functions."
}

variable logging_enabled {
  default = false
  description = "Whether cloudwatch logging should be enabled or not."
}
