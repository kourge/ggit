
all:
	go build

install:
	go install

clean:
	rm ./goit

test:
	go test ./core ./config

