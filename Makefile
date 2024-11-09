

build:
	mkdir build
	go build -o build/notifier .

test:
	go test .

clean:
	rm -rf buil
