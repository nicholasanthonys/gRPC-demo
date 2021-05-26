protoc greet/greetpb/greet.proto --go-grpc_out=.
protoc -I=. --go_out=. greet/greetpb/greet.proto