build:
	go build -o find-printers
install: build
	sudo mv find-printers /usr/local/bin/find-printers
