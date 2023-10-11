# Data-Service

The application is used to serve APIs requests for the frontend application

## Installation

Use the makefile to install the dependencies.

```bash
make init_dependency
```

Use `bash  
make copy_env` to make a local copy of the env

## Startup

Use the below command to run the application on port 8080
``bash
make run
``

## Configuration

- The service runs on port ``8080``
- The health endpoint runs on ``http://localhost:8080/api/ingest/health``
## License

[MIT](https://choosealicense.com/licenses/mit/)
