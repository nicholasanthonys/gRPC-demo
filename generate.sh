protoc greet/greetpb/greet.proto --go-grpc_out=.
protoc -I=. --go_out=. greet/greetpb/greet.proto

protoc calculator/calculatorpb/calculator.proto --go-grpc_out=.
protoc -I=. --go_out=. calculator/calculatorpb/calculator.proto