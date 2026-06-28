# 🔗 link-strat

Uma aplicação que checa códigos e status de uma lista de URLs

## Como utilizar a aplicação

> No final dos comandos, a aplicação fica disponível em `http://localhost:8080`

### 🖐 Manualmente

1. Instale o [Go](https://go.dev/dl/)
2. Execute o comando `go mod download`
3. Execute o comando `go run ./cmd/api`

### 🐳 Docker

1. Instale o [Docker](https://docs.docker.com/get-docker/)
2. Execute o comando `docker build -t link-strat .`
3. Execute o comando `docker run -d -p 8080:8080 --name link-strat-app link-strat`
   1. Se o container já existe (Já rodou o comando anteriormente), use `docker start link-strat-app`

