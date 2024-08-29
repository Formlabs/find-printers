build:
	@go build -o find-printers

install: build
	@go install .

run:
	@go run main.go