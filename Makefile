
all:
	go build

clean:
	rm ./goit

test:
	go test ./core ./config

