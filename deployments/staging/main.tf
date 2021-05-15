terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
}

provider "aws" {
  profile = "uptactics"
  region  = "us-east-1"
}

resource "aws_vpc" "vpc" {
  cidr_block = "172.30.0.0/16"

  tags = {
    Name    = "${local.service_name_env}-vpc"
  }
}

resource "aws_subnet" "private-subnet" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "172.30.1.0/24"

  tags = {
    Name    = "${local.service_name_env}-private-subnet"
  }
}

resource "aws_subnet" "public-subnet" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "172.30.2.0/24"

  tags = {
    Name    = "${local.service_name_env}-public-subnet"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name    = "${local.service_name_env}-internet-gateway"
  }
}

resource "aws_eip" "eip" {
  vpc        = true
  depends_on = [aws_internet_gateway.igw]

  tags = {
    Name    = "${local.service_name_env}-eip"
  }
}

resource "aws_nat_gateway" "ngw" {
  allocation_id = aws_eip.eip.id
  subnet_id     = aws_subnet.public-subnet.id

  tags = {
    Name    = "${local.service_name_env}-nat-gateway"
  }
}

# Public route table configuration

resource "aws_route_table" "route-table-public" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "${local.service_name_env}-route-table-public"
  }
}

resource "aws_route_table_association" "route-table-public-subnet-association" {
  subnet_id      = aws_subnet.public-subnet.id
  route_table_id = aws_route_table.route-table-public.id
}

# Private route table configuration

resource "aws_route_table" "route-table-private" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.ngw.id
  }

  tags = {
    Name = "${local.service_name_env}-route-table-private"
  }
}

resource "aws_route_table_association" "route-table-private-subnet-association" {
  subnet_id      = aws_subnet.private-subnet.id
  route_table_id = aws_route_table.route-table-private.id
}
