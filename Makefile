.PHONY: proto build test clean

all: proto build

proto:
	protoc --proto_path=proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/coordinator.proto
	protoc --proto_path=proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/storage.proto

	mv ./coordinator.pb.go coordinator/
	mv ./coordinator_grpc.pb.go coordinator/
	mv ./storage.pb.go storage/
	mv ./storage_grpc.pb.go storage/
	protoc --proto_path=proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/coordinator.proto
	protoc --proto_path=proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/storage.proto
	protoc --proto_path=proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/entities.proto
	mv ./coordinator.pb.go coordinator/
	mv ./coordinator_grpc.pb.go coordinator/
	mv ./storage.pb.go storage/
	mv ./storage_grpc.pb.go storage/
	cp ./entities.pb.go ./coordinator/entities.pb.go
	cp ./entities.pb.go ./storage/entities.pb.go
	cp ./storage/storage.pb.go ./coordinator/storage.pb.go
	cp ./storage/storage_grpc.pb.go ./coordinator/storage_grpc.pb.go

build:
	go build -o dist/storage ./storage/*.go
	go build -o dist/coordinator ./coordinator/*.go

test:
	go test ./coordinator/
	go test ./storage/

clean:
	rm -rf ./dist
	# clear protobuf generated files
	rm -rf ./*.pb.go
	rm -rf ./**/*.pb.go
