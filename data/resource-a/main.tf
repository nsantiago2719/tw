
terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.5.2"
    }
  }
}
resource "local_file" "test" {

  content  = var.content
  filename = var.filename-name
}

variable "content" {
  type = string
}

variable "filename-name" {
  type = string
}


