PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

#all: $(GO)
all: compile

compile: clean
	@protoc -I publisher/ publisher.proto --go_out=plugins=grpc:publisher

#%.pb.go: %.proto
	#@protoc --proto_path=$(GOPATH)/src:. --gogofaster_out=. $<

clean:
	@rm -f ./publisher/*.pb.go
	@rm -f ./publisher/*.go
