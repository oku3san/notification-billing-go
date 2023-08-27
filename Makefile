.PHONY: build clean

build:
		pushd ./src/notification-billing-go/ && \
    		GOOS=linux GOARCH=amd64 go build -o ./../../infra/serverless/bin/bootstrap ./main.go && \
    		popd


clean:
		@rm -rf ./src/notification-billing-go/bin/
