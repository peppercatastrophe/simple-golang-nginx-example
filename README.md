# simple-golang-nginx-example

(Work in progress)

Demonstration of Backend services using Docker, Nginx, and Go Echo

### Using Docker

```bash
docker compose up --build
```

### Using Podman

```bash
podman compose up --build
```

Or using provided Makefile if using Podman

```bash
make build
make up
```

### Swagger

```
http://localhost:8080/main/swagger/index.html
```

```
http://localhost:8080/items/swagger/index.html
```

#### Known Issues

- ~~Using `127.0.0.1` instead of `localhost` results in CORS error~~ fixed
- `http://localhost:8080/main/swagger/` (without appending `index.html`) has redirection error
- `http://localhost:8080/items/swagger/` (without appending `index.html`) has redirection error
