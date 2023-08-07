# Description
Due to implementation that relates on code generation from OpenAPI specification 
I've decided to write test task in dedicated github repository. This decision was kinda rocky 
(needed to edit some specifications in .yml ) but it works.

# How to run
- Create `config.json` file in root directory. You can copy `config_example.json` as well and rename it
- Execute `go run main.go` in root directory
- You will see info about bots in stdout

# How to regenerate api
- Install [oapi-codegen](https://github.com/deepmap/oapi-codegen) tool
- Run `make generate_api` in cmd
