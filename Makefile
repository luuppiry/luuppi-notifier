export CGO_ENABLED=0

.PHONY: build test clean run
build: clean
	mkdir -p build
	go build -o build/notifier .

test:
	go vet
	go test ./... -coverprofile cover.out

clean:
	rm -rf build

run: build
	build/notifier --configPath ./example_config.json
