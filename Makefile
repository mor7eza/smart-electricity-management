.PHONY: start
start:
	docker compose up -d --build

.PHONY: stop
stop:
	docker compose down

.PHONY: tidy
tidy:
	@ echo "Tidy -> Transmitter"
	@ cd transmitter && go mod tidy
	@ echo "Tidy -> Injestion Service"
	@ cd services/injestion-service && go mod tidy
	@ echo "Tidy -> Billing Service"
	@ cd services/billing-service && go mod tidy
