
all:
	go build

install:
	go install

clean:
	rm ./ggit

test:
	go test ./core ./config ./format ./util ./plumbing ./porcelain

