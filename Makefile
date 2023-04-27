all: build

build:
	go build -o dist/storage ./storage/*.go
	go build -o dist/coordinator ./coordinator/*.go

test:
	go test ./coordinator/
	go test ./storage/

clean:
	rm -rf ./dist