# ---- RUN & BUILD ----
run:
	go run main.go

build:
	go build -o main .

clean:
	rm -f main
	migrate-local:
	flyway -configFiles=flyway.conf migrate

# ---- FLYWAY LOCAL ----
migrate-local:
	flyway -configFiles=flyway.conf migrate

flyway-info-local:
	flyway -configFiles=flyway.conf info

reset-local:
	psql postgres://vreeni@localhost:5432/postgres -c "DROP DATABASE IF EXISTS eisbach;"
	psql postgres://vreeni@localhost:5432/postgres -c "CREATE DATABASE eisbach OWNER vreeni;"
	make migrate-local

# ---- FLYWAY PROD (Neon) ----
migrate-prod:
	flyway -configFiles=flyway-prod.conf migrate

flyway-info-prod:
	flyway -configFiles=flyway-prod.conf info