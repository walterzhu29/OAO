export GOBIN=$(go env GOPATH)/bin                             
export PATH=$PATH:$GOBIN


go mod init submitter 


go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


protoc --go_out=. --go-grpc_out=. model_server.proto