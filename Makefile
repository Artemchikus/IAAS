build:
	@go build -o bin/IAAS cmd/IAAS/main.go

run: build
	@./bin/IAAS

test:
	@go test -v ./...

seed: build
	@./bin/IAAS --seed

rebuildpq:
	@docker stop some-postgres
	@docker rm some-postgres
	@docker run --name some-postgres -e POSTGRES_PASSWORD=iaas -e POSTGRES_DB=iaas -p 5432:5432 -d postgres

buildpq:
	@docker run --name some-postgres -e POSTGRES_PASSWORD=iaas -e POSTGRES_DB=iaas -p 5432:5432 -d postgres

buildtestpq:
	@docker run --name some-test-postgres -e POSTGRES_PASSWORD=iaas -e POSTGRES_DB=iaas-test -p 5433:5432 -d postgres
	
rebuildtestpq:
	@docker stop some-test-postgres
	@docker rm some-test-postgres
	@docker run --name some-test-postgres -e POSTGRES_PASSWORD=iaas -e POSTGRES_DB=iaas-test -p 5433:5432 -d postgres