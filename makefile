.PHONY: generate_api
generate_api:
	openapi-generator generate -i internal/api/jaicp/jaicp_specification.yml -g go -o internal/api/jaicp/generated --git-host github.com --git-user-id jadegopher  --git-repo-id task