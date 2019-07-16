AWS_REGION=eu-central-1

.PHONY: deps clean build

deps:
	dep ensure -v
#	go get -u ./...

clean: 
	rm -rf ./build/bin/*
	rm -rf ./build/archive/*

build:
	GOOS=linux GOARCH=amd64 go build -o build/bin/sam-init-template-runtime-go ./cmd/sam-init-template-runtime-go
	zip -j build/archive/sam-init-template-runtime-go.zip build/bin/sam-init-template-runtime-go

debug:
	sam local invoke -e event.json HelloWorldFunction --debug

publish:
	aws lambda update-function-code --function-name sam-init-template-runtime-go --zip-file fileb://build/archive/sam-init-template-runtime-go.zip
#	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket scherkesov1006-lambda-deploy
#	sam publish --template packaged.yaml --region ${AWS_REGION}