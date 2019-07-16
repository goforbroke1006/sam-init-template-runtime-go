include .env
export $(shell sed 's/=.*//' .env)

.PHONY: deps clean build

deps:
	dep ensure -v

clean: 
	rm -rf ./build/bin/*
	rm -rf ./build/archive/*

build:
	GOOS=linux GOARCH=amd64 go build -o build/bin/sam-init-template-runtime-go ./cmd/sam-init-template-runtime-go

build-debug:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/sam-init-template-runtime-go ./cmd/sam-init-template-runtime-go

start-api:
	go get -u github.com/go-delve/delve/cmd/dlv
	go build -o ./build/debugger/dlv github.com/go-delve/delve/cmd/dlv
	sam local start-api --port 3000 --debug-port 5986 --debugger-path ./build/debugger/ --debug-args '-delveAPI=2'

publish:
	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket scherkesov1006-lambda-deploy
	sam deploy --template-file packaged.yaml --stack-name sam-init-template-runtime-go --capabilities CAPABILITY_IAM
