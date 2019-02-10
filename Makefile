.PHONY: all build

all: build

get:
	@go get || true

build: get
	@go build -i -v -o pssh ./

clean:
	@rm ./pssh
