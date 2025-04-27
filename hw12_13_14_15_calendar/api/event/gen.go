//go:generate protoc --proto_path=../../api --go_out=../../pb/event --go-grpc_out=../../pb/event --experimental_allow_proto3_optional event/event.proto
package event
