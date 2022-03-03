APP := gitmono
GO_COVER_FILE ?= "coverage.out"

build:
	CGO_ENABLED=0 go build -o $(APP)/$(APP)  cmd/$(APP)/*.go

test:
	go test ./...

test-cover:
	[ ! -e $(GO_COVER_FILE) ] || rm $(GO_COVER_FILE)
	go test ./... --count=1 -race -covermode=atomic -coverprofile=$(GO_COVER_FILE)
	go tool cover -func $(GO_COVER_FILE)
