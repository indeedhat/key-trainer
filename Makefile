.PHONY build
build:
	go build -o ./ ./...
	
run: build
	./key-trainer
