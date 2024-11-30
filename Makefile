
.PHONY: build test clean run
build: clean
	-mkdir build
	go build -o build/notifier .

test:
	go test .

clean:
	rm -rf build

run: build
	build/notifier --configPath ./example_config.json
