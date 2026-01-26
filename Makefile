
.PHONY: all database kafka server

all: database kafka server

database:
	@echo "Running Database"
	cd database && air

kafka:
	@echo "Running Kafka"
	cd kafka && air

server:
	@echo "Starting Server"
	cd server && air

wait:
	@wait

