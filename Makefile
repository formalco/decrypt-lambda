.PHONY: build clean deploy

build:
	go mod tidy
	env CGOENABLED=0 GOARCH=arm64 GOOS=linux go build -o bootstrap .
clean:
	rm -rf ./bin bootstrap

deploy-sls: clean build
	cd serverless && sls deploy --verbose

deploy-terraform: clean build
	zip bootstrap.zip bootstrap
	cd terraform && terraform init && terraform apply