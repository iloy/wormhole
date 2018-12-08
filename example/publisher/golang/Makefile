
all: example_publisher_go

example_publisher_go: main.go ../../../lib/lib.go ../../../protocol/pubsub.pb.go
	go build -race -o $@

test: example_publisher_go
	./$<

.PHONY: clean
clean:
	rm -f example_publisher_go

