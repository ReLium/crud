# crud
CRUD demo project in Golang

## Usage

- `make build`
- `docker-compose up -d`
- `make generate`
- go to http://127.0.0.1:8080/
- swagger UI http://127.0.0.1:8082/?url=http://127.0.0.1:8080/swagger.yaml


## TODOs
- refine errors handling
- refine mongodb context usage
- another abstraction layer to handle another entity types
- more unit tests
- implement zerolog for json logging (?)
