include .env
export $(shell sed 's/=.*//' .env)

.PHONY: deps clean build

deps:
	dep ensure -v

clean: 
	rm -rf ./build/bin/*
	rm -rf ./build/archive/*

build:
	GOOS=linux GOARCH=amd64 go build -o build/bin/sam-init-template-runtime-go  ./cmd/sam-init-template-runtime-go
	GOOS=linux GOARCH=amd64 go build -o build/bin/current-time                  ./cmd/current-time
	GOOS=linux GOARCH=amd64 go build -o build/bin/hello-username                ./cmd/hello-username

build-debug:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/sam-init-template-runtime-go ./cmd/sam-init-template-runtime-go
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/current-time                 ./cmd/current-time
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/hello-username               ./cmd/hello-username

start-api:
	go get -u github.com/go-delve/delve/cmd/dlv
	go build -o ./build/debugger/dlv github.com/go-delve/delve/cmd/dlv
	sam local start-api --port 3000 --debug-port 5986 --debugger-path ./build/debugger/ --debug-args '-delveAPI=2'

publish:
	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket ${AWS_DEPLOY_BUCKET_NAME}
	sam deploy --template-file packaged.yaml --stack-name sam-init-template-runtime-go --capabilities CAPABILITY_IAM


debug/CurrentTimeFunction:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/current-time ./cmd/current-time
	sam local invoke --event cmd/current-time/event.json CurrentTimeFunction --debug

debug/HelloUsernameFunction:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/hello-username ./cmd/hello-username
	sam local invoke --event cmd/hello-username/event.json HelloUsernameFunction --debug

