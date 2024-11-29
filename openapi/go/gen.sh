#!/bin/sh

# NOTE: This script assumes it is run from project root

mkdir -p openapi/go/bureaucrat

go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 -config openapi/go/oapi-codegen.yml openapi/openapi.yml

go mod tidy

go run github.com/vektra/mockery/v2@v2.42.0