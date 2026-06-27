FROM golang:1.26.4
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /link-strat ./cmd/api
EXPOSE 8080
CMD ["/link-strat"]