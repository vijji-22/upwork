# Temple Management System

## Table of Contents
- [Service Description](#service-description)
- [REST Endpoints](#rest-endpoints)
- [Installation](#installation)
- [Start Server](#start-server)
  - [Using Docker](#using-docker)
  - [Migration](#migration)
  - [Debugging with VSCode](#debugging-with-vscode)
- [Contributing](#contributing)
- [Tests](#tests)
- [License](#license)

## Service Description
Replace this section with a brief description of what your project does and its purpose.

## REST Endpoints
## REST Endpoints

### Profiling Endpoints
- **GET** /debug/pprof/
- **GET** /debug/pprof/cmdline
- **GET** /debug/pprof/profile
- **POST** /debug/pprof/symbol
- **GET** /debug/pprof/symbol
- **GET** /debug/pprof/trace
- **GET** /debug/pprof/allocs
- **GET** /debug/pprof/block
- **GET** /debug/pprof/goroutine
- **GET** /debug/pprof/heap
- **GET** /debug/pprof/mutex
- **GET** /debug/pprof/threadcreate

### Template Endpoints
- **GET** /api/v1/template/
- **GET** /api/v1/template/:id
- **POST** /api/v1/template/
- **PUT** /api/v1/template/:id
- **DELETE** /api/v1/template/:id

### Temple Endpoints
- **GET** /api/v1/temple/
- **GET** /api/v1/temple/:id
- **POST** /api/v1/temple/
- **PUT** /api/v1/temple/:id
- **DELETE** /api/v1/temple/:id

### Temple Config Meta Endpoints
- **GET** /api/v1/temple_config_meta/
- **GET** /api/v1/temple_config_meta/:id
- **POST** /api/v1/temple_config_meta/
- **PUT** /api/v1/temple_config_meta/:id
- **DELETE** /api/v1/temple_config_meta/:id
- **GET** /api/v1/temple_config_meta/tree

### Temple Config Value Endpoints
- **GET** /api/v1/temple_config_value/
- **GET** /api/v1/temple_config_value/:id
- **POST** /api/v1/temple_config_value/
- **PUT** /api/v1/temple_config_value/:id
- **DELETE** /api/v1/temple_config_value/:id
- **POST** /api/v1/temple_config_value/tree/init/:temple_id
- **GET** /api/v1/temple_config_value/tree/:temple_id

### Access Endpoints
- **GET** /api/v1/access/
- **GET** /api/v1/access/:id
- **POST** /api/v1/access/
- **PUT** /api/v1/access/:id
- **DELETE** /api/v1/access/:id

### Role Endpoints
- **GET** /api/v1/role/
- **GET** /api/v1/role/:id
- **POST** /api/v1/role/
- **PUT** /api/v1/role/:id
- **DELETE** /api/v1/role/:id

### Role Access Mapping Endpoints
- **GET** /api/v1/role_access_mapping/
- **GET** /api/v1/role_access_mapping/:id
- **POST** /api/v1/role_access_mapping/
- **PUT** /api/v1/role_access_mapping/:id
- **DELETE** /api/v1/role_access_mapping/:id

### Authentication Endpoints
- **POST** /api/v1/auth/login

### User Endpoints
- **GET** /api/v1/user/
- **GET** /api/v1/user/:id
- **POST** /api/v1/user/
- **PUT** /api/v1/user/:id
- **DELETE** /api/v1/user/:id

### User Role Mapping Endpoints
- **GET** /api/v1/user_role_mapping/
- **GET** /api/v1/user_role_mapping/:id
- **POST** /api/v1/user_role_mapping/
- **PUT** /api/v1/user_role_mapping/:id
- **DELETE** /api/v1/user_role_mapping/:id
- **GET** /api/v1/user_role_mapping/detail
- **GET** /api/v1/user_role_mapping/detail/:user_id


## Installation
```bash
make build
```

## Start Server
### Using Docker
If we want to run DB and App both in docker then `docker-compose up --build`
If we want to run DB only `docker-compose up -d db`
If we already have Database, just update connection details in `docker-compose.yml` update ENVs for `services > app > environment` and run `docker-compose up --build app`
```yaml
    environment:
      - DB_HOST=db
      - DB_USER=pradip
      - DB_PASSWORD=password
      - DB_DATABASE=test
```
This process will start the server on port `8080`. If you want to run on any other server then pass `APP_PORT` in `docker-compose app` e.g `APP_PORT=8081 docker-compose app`

### Migration
- using github.com/golang-migrate/migrate for migration
- need to install migration => brew install golang-migrate
- create a new migration file `migrate create -ext sql -dir migrations -seq create_access_table`

### Debugging with VSCode
.vscode/launch.json
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/server/main.go",
            "args": ["runserver"],
            "env": {
                "DB_HOST": "127.0.0.1",
                "DB_USER":"pradip",
                "DB_PASSWORD":"password",
                "DB_DATABASE":"test",
                "DB_PORT":"54321",
                "MIGRATION_FILE_PATH": "${workspaceFolder}/migrations",
            }
        }
    ]
}
```

```

## Contributing
Replace this section with guidelines for contributors.

## Tests
Replace this section with information on how to run tests.

## License
Replace this section with information about your project's license.
```

Feel free to customize the content based on your project specifics.