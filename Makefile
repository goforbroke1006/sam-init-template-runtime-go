AWS_CF_STACK=sam-init-template-runtime-go

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
	GOOS=linux GOARCH=amd64 go build -o build/bin/retrieve-companies-data       ./cmd/retrieve-companies-data

build-debug:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/sam-init-template-runtime-go ./cmd/sam-init-template-runtime-go
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/current-time                 ./cmd/current-time
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/hello-username               ./cmd/hello-username
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/retrieve-companies-data      ./cmd/retrieve-companies-data

start-api:
	go get -u github.com/go-delve/delve/cmd/dlv
	go build -o ./build/debugger/dlv github.com/go-delve/delve/cmd/dlv
	sam local start-api --port 3000 --debug-port 5986 --debugger-path ./build/debugger/ --debug-args '-delveAPI=2'

publish:
	aws s3 mb s3://${AWS_DEPLOY_BUCKET_NAME} || true
	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket ${AWS_DEPLOY_BUCKET_NAME}
	sam deploy --template-file packaged.yaml --stack-name ${AWS_CF_STACK} --capabilities CAPABILITY_IAM

logs/publish:
	aws cloudformation describe-stack-events --stack-name sam-init-template-runtime-go --max-items 25




debug/current-time:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/current-time ./cmd/current-time
	sam local invoke --event cmd/current-time/event.json CurrentTimeFunction --debug

debug/hello-username:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/hello-username ./cmd/hello-username
	sam local invoke --event cmd/hello-username/event.json HelloUsernameFunction --debug

debug/retrieve-companies-data:
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/retrieve-companies-data ./cmd/retrieve-companies-data
	sam local generate-event s3 put --bucket test-bucket --debug | sam local invoke RetrieveCompaniesDataFunction


logs/retrieve-companies-data:
	sam logs -n RetrieveCompaniesDataFunction --stack-name ${AWS_CF_STACK} --tail
