terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }

  backend "s3" {
    region         = "ap-northeast-1"
    bucket         = "oku3san-tfstate-bucket"
    key            = "notification-billing-go-tfstate"
    role_arn       = "arn:aws:iam::950007738962:role/terraform"
  }
}

provider "aws" {
  region = "ap-northeast-1"
  assume_role {
    role_arn = "arn:aws:iam::950007738962:role/terraform"
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
