.PHONY: build clean deploy

build:
	go mod tidy
	env CGOENABLED=0 GOARCH=arm64 GOOS=linux go build -o bootstrap ./main.go
clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose
