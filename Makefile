APP := gitmono

build:
	CGO_ENABLED=0 go build -o $(APP)/$(APP)  cmd/$(APP)/*.go
