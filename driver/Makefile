APP_NAME := driver

# Default target
build: clean depend go-build

clean:
	rm -f ./$(APP_NAME)

depend:
	go get -u ./...

go-build:
	go build

# None of the Make tasks generate files with the name of the task, so all must be declared as 'PHONY'
.PHONY: clean depend build go-build
