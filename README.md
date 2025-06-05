![Build Status](https://github.com/Gburan/pvz-service/actions/workflows/build.yml/badge.svg?job=build)
![Tests Status](https://github.com/Gburan/pvz-service/actions/workflows/test.yml/badge.svg?job=test)
![Lint Status](https://github.com/Gburan/pvz-service/actions/workflows/lint.yml/badge.svg?job=lint)
[![codecov](https://codecov.io/gh/Gburan/pvz-service/graph/badge.svg?token=8OFOUGXOYL)](https://codecov.io/gh/Gburan/pvz-service)

# PVZ Service

## 1. Description

A service for working with PVZ. Full text of the assignment can be found [here](./task).
Additionally, done: {`docs generating`,`more integrational tests`,`CI`}.

1. The `/dummyLogin` method: Accepts the user type (employee, moderator) and returns a token with the appropriate
   access level (employee or moderator).
2. The `/register` method: Accepts user registration data (mail, password, and type — employee or moderator).
   Creates a user and returns a token for authorization.
3. The `/login` method: Accepts mail and password, returns a token for an authorized user with the appropriate
   access level.
4. The `/pvz` method: For moderators only. Accepts data for the creation of a new PVZ in one of the three
   cities (Moscow, St. Petersburg, Kazan). Returns full information about the created IDP or an error if the city is not
   fits.
5. The `/receptions` method: Only for authorized users with the role of "employee". Initializes
   acceptance of goods in the PVZ. Returns the result of the operation or an error if the previous reception was not
   closed.
6. The `/products` method: Only for authorized users with the role of "employee". Adds products to
   the current acceptance of goods. Returns an error if there is no pending acceptance.
7. The `/pvz/{pvzId}/delete_last_product` method: Only for authorized users with the role of "employee". Removes
   products from
   open acceptance according to the LIFO principle (the last added product is removed first).
8. The `/pvz/{pvzId}/close_last_reception` method: Only for authorized users with the role of "employee". Closes the
   current acceptance
   products. Returns an error if acceptance has already been closed.
9. The `/pvz` method: Only for authorized users with the role of "employee" or "moderator". Returns
   a list of PVZ with information about the acceptance of goods for a specified period of time, with pagination support.

## 2. Configuration

| Name                   | Type   | Default value                                                                      | Description                                               |
|------------------------|--------|------------------------------------------------------------------------------------|-----------------------------------------------------------|
| REST_ADDRESS           | String | `:8080`                                                                            | REST server address                                       |
| GRPC_ADDRESS           | String | `:3000`                                                                            | gRPC server address                                       |
| POSTGRES_CONN          | String |                                                                                    | PostgreSQL connection string                              |
| REST_CONN_SETTINGS     | String | `read_timeout: 5s, write_timeout: 5s, idle_timeout: 5m`                            | REST connection settings                                  |
| GRPC_CONN_SETTINGS     | String | `max_conn_idle: 5m, max_conn_age: 10m`                                             | gRPC connection settings                                  |
| POSTGRES_POOL_SETTINGS | String | `max_conns: 15, min_idle_conns: 5, max_conn_idle_time: 5m, max_conn_lifetime: 10m` | PostgreSQL connection pool settings                       |
| MIGRATIONS_DIR         | String | `./migrations`                                                                     | Directory for database migrations                         |
| ALLOWED_CATEGORIES     | List   | `["электроника", "одежда", "обувь"]`                                               | List of allowed product categories                        |
| ALLOWED_CITIES         | List   | `["Москва", "Санкт-Петербург", "Казань"]`                                          | List of allowed cities                                    |
| ALLOWED_USERS          | List   | `["moderator", "employee"]`                                                        | List of allowed user roles                                |
| JWT_TOKEN              | String |                                                                                    | Predefined JWT token for authentication                   |
| LOGGING_OUTPUT         | String | `"stdout"`                                                                         | Log output destination ("stdout", "stderr", or file path) |
| LOGGING_LEVEL          | String | `"debug"`                                                                          | Log level ("debug", "info", "warn", "error")              |

## 3. How to run

```
go mod tidy
go generate ./...
docker-compose up -d
```

or (first time)

```
make run-compose-f
```

or (with container rebuild)

```
make run-compose-b
```

or

```
make run-compose
```

## 4. Tests & Code generation

### 4.1 Generate mocks & gRPC & DTOs for Endpoints

```
go generate ./...
```

or

```
make gen
```

### 4.2 Run tests

```
go clean -testcache
go test ./...
```

or

```
make run-test-clean
```

### 4.3 Postman

`For REST testing import the config` [file](./test/postman/PVZ.postman_collection.json)

`For gRPC testing import the .proto` [file](./api/v1/proto/pvz.proto)

...into the postman for manual testing.

## 5. Metrics

Metrics are available [here](http://localhost:3030/). [Login / Pass]: [admin / admin]

1)Move to `Connections/Data sources`
-> `Add data source` -> `Prometheus` -> enter `http://prometheus:9090` into `Prometheus server URL` field
-> `Save & Test`.

2)Move to `Dashboard` -> `New` -> `Import` -> import the [file](./metrics/grafana/pvz_service-1748169298112.json)


## 6. Docs

### 6.1 gRPC Docs

Docs is [here](./docs/grpc/grpc.md)

### 6.2 Rest Docs

Swagger is available [here](http://localhost:8080/swagger/)

Regenerate doc: `swag init -g /internal/app/setup.go --output docs/rest`
