package proto

//go:generate protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. racing/racing.proto --experimental_allow_proto3_optional
