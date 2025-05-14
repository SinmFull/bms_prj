include .envrc

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: run/api
run/api:
	go run ./cmd/api  -db-dsn=${BMS_DB_DSN}

.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${BMS_DB_DSN} up

.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running up migrations down...'
	migrate -path ./migrations -database ${BMS_DB_DSN} down

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #
production_host_ip = '206.189.92.86'
## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh bms@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -rP --delete ./bin/linux_amd64/api ./migrations bms@${production_host_ip}:~
	ssh -t bms@${production_host_ip} 'migrate -path ~/migrations -database ${BMS_DB_DSN} up'