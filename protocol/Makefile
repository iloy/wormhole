
#all: pubsub.pb.go pubsub_pb2.py pubsub_pb2_grpc.py Pubsub.java pubsub_pb.js pubsub_grpc_web_pb.js pubsub.pb.cc pubsub.pb.h
all: pubsub.pb.go pubsub_pb2.py pubsub_pb2_grpc.py

pubsub.pb.go: pubsub.proto
	go get -u github.com/gogo/protobuf/protoc-gen-gogofaster
	protoc --gogofaster_out=plugins=grpc:. pubsub.proto

pubsub_pb2.py: pubsub.proto
	pip install -U grpcio-tools
	pip install -U googleapis-common-protos
	python -m grpc_tools.protoc -I. -I../vendor --python_out=. --grpc_python_out=. pubsub.proto
pubsub_pb2_grpc.py: pubsub.proto
	pip install -U grpcio-tools
	pip install -U googleapis-common-protos
	python -m grpc_tools.protoc -I. -I../vendor --python_out=. --grpc_python_out=. pubsub.proto

# FIXME!
#Pubsub.java: pubsub.proto
#	protoc -I=. --java_out=.. pubsub.proto

protoc-gen-grpc-web:
	wget https://github.com/grpc/grpc-web/releases/download/1.0.3/protoc-gen-grpc-web-1.0.3-linux-x86_64
	mv protoc-gen-grpc-web-1.0.3-linux-x86_64 protoc-gen-grpc-web
	chmod +x protoc-gen-grpc-web
	mv ~/local/bin

pubsub_pb.js: pubsub.proto
	protoc -I=. --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:. pubsub.proto
pubsub_grpc_web_pb.js: pubsub.proto
	protoc -I=. --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:. pubsub.proto

# FIXME!
#pubsub.pb.cc: pubsub.proto
#	protoc -I=. --cpp_out=. pubsub.proto
#pubsub.pb.h: pubsub.proto
#	protoc -I=. --cpp_out=. pubsub.proto

clean:
	#rm -f *.go *.py *.java *.js *.cc *.h

