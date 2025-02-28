.PHONY: cdown cbuild cup clog restart install prune

cdown:
	docker-compose down

cbuild:
	docker-compose build

cup:
	docker-compose up -d

clog:
	docker ps && docker-compose logs -f exchanger

restart: cdown cbuild cup

install: cbuild cup clog

prune:
	docker-compose down -v --rmi all && docker system prune -f

test:
	grpcurl -plaintext -d '{}' localhost:50051 exchange.ExchangeService/GetExchangeRates