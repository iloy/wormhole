
all: wormhole

SRCS := main.go
SRCS += version.go
SRCS += beancounter/*.go
SRCS += protocol/pubsub.pb.go
SRCS += publisher/*.go 
SRCS += seqno/*.go 
SRCS += server/*.go
SRCS += subscriber/*.go 
SRCS += token/*.go
SRCS += topic/*.go
SRCS += topicmanager/*.go


BUILDER_NAME = wormhole_builder


wormhole: $(SRCS)
	go build -a -ldflags "-X main.revision=`./revision.sh` -X main.modified=`./modified.sh` -X main.builddatetime=`./built.sh`"
	go install ./...

.PHONY: race
race:
	go build -a -race -gcflags='all=-N -l' -ldflags "-X main.revision=`./revision.sh` -X main.modified=`./modified.sh` -X main.builddatetime=`./built.sh`"

.PHONY: build
build:
	make clean
	-docker rm -f $(BUILDER_NAME)
	docker run -d -it --name $(BUILDER_NAME) -v `pwd`:/go/src/github.com/iloy/wormhole golang:1.11 /bin/bash
	docker exec -t $(BUILDER_NAME) /bin/bash -c "go get -u github.com/kardianos/govendor && cd /go/src/github.com/iloy/wormhole && govendor sync"
	docker exec -t $(BUILDER_NAME) /bin/bash -c "cd /go/src/github.com/iloy/wormhole && make"
	-docker rm -f $(BUILDER_NAME)

protocol/pubsub.pb.go: protocol/pubsub.proto
	make -C protocol


localhost: wormhole
	./wormhole


test: race
	go test -v ./...
	./wormhole --version
	./wormhole

.PHONY: clean
clean:
	rm -f wormhole
	make -C protocol clean

