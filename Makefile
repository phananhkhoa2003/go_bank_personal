postgres:
	docker run --name postgres17.5 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.5-alpine
simple_bank:
	docker run --name simplebank --network bank-network  -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres17.5:5432/simple_bank?sslmode=disable" simplebank:latest
createdb:
	docker exec -it postgres17.5 createdb --username=root --owner=root simple_bank

dropdb:
	
	docker exec -it postgres17.5 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateupaws:
	migrate -path db/migration -database "postgresql://root:i7mEA2OSZNQKTdXrVOMP@simple-bank.chu2eo4g2mkm.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migrateupaws1:
	migrate -path db/migration -database "postgresql://root:i7mEA2OSZNQKTdXrVOMP@simple-bank.chu2eo4g2mkm.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up 1
migratedownaws:
	migrate -path db/migration -database "postgresql://root:i7mEA2OSZNQKTdXrVOMP@simple-bank.chu2eo4g2mkm.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose down
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
migratedownaws1:
	migrate -path db/migration -database "postgresql://root:i7mEA2OSZNQKTdXrVOMP@simple-bank.chu2eo4g2mkm.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

# Benchmark commands
benchmark:
	go test -bench=. -benchmem ./...

benchmark-api:
	go test -bench=. -benchmem ./api

benchmark-db:
	go test -bench=. -benchmem ./db/sqlc

benchmark-token:
	go test -bench=. -benchmem ./token

benchmark-util:
	go test -bench=. -benchmem ./util

benchmark-report:
	go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...

# Run comprehensive benchmark tests and save results to benchmark_results folder
benchmark-save:
	powershell -Command "Write-Host 'Running comprehensive benchmark tests...'"
	powershell -Command "if (!(Test-Path 'benchmark_results')) { New-Item -ItemType Directory -Path 'benchmark_results' }"
	powershell -Command "Write-Host 'Testing database layer...'; go test -bench=Benchmark -benchmem -v ./db/sqlc | Tee-Object -FilePath \"benchmark_results/benchmark_results_sqlc_$$(Get-Date -Format 'yyyyMMdd_HHmmss').txt\""
	powershell -Command "Write-Host 'Testing API layer...'; go test -bench=Benchmark -benchmem -v ./api | Tee-Object -FilePath \"benchmark_results/benchmark_results_api_$$(Get-Date -Format 'yyyyMMdd_HHmmss').txt\""
	powershell -Command "Write-Host 'Testing token layer...'; go test -bench=Benchmark -benchmem -v ./token | Tee-Object -FilePath \"benchmark_results/benchmark_results_token_$$(Get-Date -Format 'yyyyMMdd_HHmmss').txt\""
	powershell -Command "Write-Host 'Testing util layer...'; go test -bench=Benchmark -benchmem -v ./util | Tee-Object -FilePath \"benchmark_results/benchmark_results_util_$$(Get-Date -Format 'yyyyMMdd_HHmmss').txt\""
	powershell -Command "Write-Host 'Benchmark tests completed! Results saved in benchmark_results/ folder'"

loadtest:
	cd loadtest && go run main.go

# Mock database server for load testing
mock-server:
	cd mock_server && go run main.go

# Run load test against mock server (recommended for performance testing)
loadtest-mock:
	cd loadtest_mock && go run main.go

# Run load test against real server (requires server to be running on 8080)
loadtest-real:
	cd loadtest && go run main.go

server: 
	go run main.go

mock:
	mockgen -package mockdb  -destination db/mock/store.go simple_bank/db/sqlc Store 

# Display help for benchmark commands
benchmark-help:
	@echo "Available benchmark commands:"
	@echo "  make benchmark           - Run all benchmarks"
	@echo "  make benchmark-api       - Run API layer benchmarks"
	@echo "  make benchmark-db        - Run database layer benchmarks"
	@echo "  make benchmark-token     - Run token layer benchmarks"
	@echo "  make benchmark-util      - Run utility layer benchmarks"
	@echo "  make benchmark-report    - Run benchmarks with CPU/memory profiling"
	@echo "  make benchmark-save      - Run comprehensive benchmarks and save to benchmark_results/"
	@echo ""
	@echo "Results are saved in the benchmark_results/ folder with timestamps."

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server test mock migrateup1 migratedown1 simple_bank migrateupaws migrateupaws1 migratedownaws migratedownaws1 benchmark benchmark-api benchmark-db benchmark-token benchmark-util benchmark-report benchmark-save benchmark-help loadtest mock-server loadtest-mock loadtest-real