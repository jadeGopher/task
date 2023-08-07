.PHONY: generate_api
generate_api:
	oapi-codegen -package api internal/api/jaicp_specification_edited.yml > internal/api/jaicp.gen.go