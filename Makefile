
.PHONY: build test clean
build: clean
	-mkdir build
	go build -o build/notifier .

test:
	go test .

clean:
	rm -rf build
