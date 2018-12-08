
all: example_subscriber_go

example_subscriber_go: main.go ../../../lib/lib.go ../../../protocol/pubsub.pb.go
	go build -race -o $@

test: example_subscriber_go
	./$<

.PHONY: clean
clean:
	rm -f example_subscriber_go

