> Dockerfile has shared network for the same subnet, both should be on the same localhost

Requirements:
1) Docker
2) Docker-compose
3) Make
4) Ubuntu
5) Go (Mainly for grpcurl test)

# Default .env file:

- EXCHANGER_PORT=50051
- DB_URL=postgres://postgres:postgres@exchanger-db:5432/exchange?sslmode=disable
-  POSTGRES_USER=postgres
-  POSTGRES_PASSWORD=postgres
-  POSTGRES_DB=exchange

> Makefile scripts for installations
> - make install - clean installation
> - make clog - docker logs
> - make restart - down, build, up for docker-compose
> - make test - grpc testing, need to install grpcurl ```go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest```
