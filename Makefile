deps:
	go get -u github.com/nats-io/nats
	go get -u github.com/supu-io/payload
	go get -u github.com/supu-io/messages
dev-deps:
	go get -u github.com/gorilla/mux
	go get -u github.com/smartystreets/goconvey/convey
test: 
	go test
lint:
	golint
