terraform {
  required_providers {
    vultr = {
      source  = "vultr/vultr"
      version = "2.1.2"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 1.4.0"
	  }
  }
  required_version = "~> 0.14.0"
}
