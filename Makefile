app-serve:
	docker-compose up --build -d

.PHONY: app-serve

run-tests:
	go test app/main_test.go -test.timeout=0

.PHONY: run-tests

app-down:
	docker-compose down

.PHONY: app-down