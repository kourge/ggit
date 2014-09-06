
all:
	go build

install:
	go install

clean:
	rm ./ggit

test:
	go test ./...

