.PHONY: openapi-iam

iam-openapi:
	go generate services/iam/api/api.gen.go

setup:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
