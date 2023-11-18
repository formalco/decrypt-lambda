.PHONY: build clean deploy

build:
	go mod tidy
	env CGOENABLED=0 GOARCH=arm64 GOOS=linux go build -o bootstrap .
clean:
	rm -rf ./bin bootstrap

deploy: clean build
	sls deploy --verbose
