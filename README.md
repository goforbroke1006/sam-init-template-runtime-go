# sam-init-template-runtime-go

This is a sample template for SAM project (runtime=go1.x)

Below is a brief explanation of what we have generated for you:

```bash
.
├── Makefile                            <-- Make to automate build
├── README.md                           <-- This instructions file
├── .env                                <-- Var export for Makefile (git ignored)
├── .env.dist                           <-- Sample of ".env"
├── cmd
│   └── current-time                    <-- Root of runner for lambda function
│   │   ├── main.go                     <-- Lambda function code
│   │   └── main_test.go                <-- Unit tests
│   └── hello-username                  <-- Root of runner for lambda function
│   │   ├── event.json                  <-- Debug request
│   │   ├── main.go                     <-- Lambda function code
│   │   └── main_test.go                <-- Unit tests
│   └── sam-init-template-runtime-go    <-- Root of runner for lambda function
│       ├── main.go                     <-- Lambda function code
│       └── main_test.go                <-- Unit tests
└── template.yaml                       <-- SAM-friendly configuration
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* [AWS SAM CLI](https://docs.aws.amazon.com/en_us/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies

```shell
make deps
```

### Building

```shell
make build
```

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://localhost:3000/hello`

**SAM CLI** is used to emulate both Lambda and API Gateway locally and uses our `template.yaml` to understand how to bootstrap this environment (runtime, where the source code is, etc.) - The following excerpt is what the CLI will read in order to initialize an API and its routes:

```yaml
...
Events:
    HelloWorld:
        Type: Api # More info about API Event Source: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
        Properties:
            Path: /hello
            Method: get
```

## Packaging and deployment

AWS Lambda Python runtime requires a flat folder with all dependencies including the application. SAM will use `CodeUri` property to know where to look up for both application and dependencies:

```yaml
...
    FirstFunction:
        Type: AWS::Serverless::Function
        Properties:
            CodeUri: hello_world/
            ...
```

First and foremost, we need a `S3 bucket` where we can upload our Lambda functions packaged as ZIP before we deploy anything - If you don't have a S3 bucket to store code artifacts then this is a good time to create one:

```bash
aws s3 mb s3://BUCKET_NAME
```

Next, run the following command to package our Lambda function to S3:

```bash
sam package \
    --output-template-file packaged.yaml \
    --s3-bucket REPLACE_THIS_WITH_YOUR_S3_BUCKET_NAME
```

Next, the following command will create a Cloudformation Stack and deploy your SAM resources.

```bash
sam deploy \
    --template-file packaged.yaml \
    --stack-name remove-me-1 \
    --capabilities CAPABILITY_IAM
```

> **See [Serverless Application Model (SAM) HOWTO Guide](https://github.com/awslabs/serverless-application-model/blob/master/HOWTO.md) for more details in how to get started.**

After deployment is complete you can run the following command to retrieve the API Gateway Endpoint URL:

```bash
aws cloudformation describe-stacks \
    --stack-name remove-me-1 \
    --query 'Stacks[].Outputs'
``` 

### Testing

We use `testing` package that is built-in in Golang and you can simply run the following command to run our tests:

```shell
go test -v ./hello-world/
```

## AWS CLI commands

AWS CLI commands to package, deploy and describe outputs defined within the cloudformation stack:

```bash
sam package \
    --template-file template.yaml \
    --output-template-file packaged.yaml \
    --s3-bucket REPLACE_THIS_WITH_YOUR_S3_BUCKET_NAME

sam deploy \
    --template-file packaged.yaml \
    --stack-name remove-me-1 \
    --capabilities CAPABILITY_IAM \
    --parameter-overrides MyParameterSample=MySampleValue

aws cloudformation describe-stacks \
    --stack-name remove-me-1 --query 'Stacks[].Outputs'
```

## Bringing to the next level

Here are a few ideas that you can use to get more acquainted as to how this overall process works:

* Create an additional API resource (e.g. /hello/{proxy+}) and return the name requested through this new path
* Update unit test to capture that
* Package & Deploy

Next, you can use the following resources to know more about beyond hello world samples and how others structure their Serverless applications:

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)

## Useful links

* [Index of Lambda apps](https://docs.aws.amazon.com/en_us/lambda/latest/dg/deploying-lambda-apps.html)
* [SAM CLI quick start](https://docs.aws.amazon.com/en_us/serverless-application-model/latest/developerguide/serverless-quick-start.html)
* [SAM CLI deploy](https://docs.aws.amazon.com/en_us/serverless-application-model/latest/developerguide/serverless-deploying.html)
* Example [template.yaml](https://docs.aws.amazon.com/en_us/lambda/latest/dg/with-s3-example-use-app-spec.html)
* Resources types in template.yaml [here](https://docs.aws.amazon.com/en_us/serverless-application-model/latest/developerguide/serverless-sam-template.html)
* Golang samples for S3 API [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/)
* Here [SAM for scheduled lambda function](https://docs.aws.amazon.com/en_us/lambda/latest/dg/with-scheduledevents-example-use-app-spec.html) and [CRON annotations](https://docs.aws.amazon.com/en_us/lambda/latest/dg/tutorial-scheduled-events-schedule-expressions.html)
