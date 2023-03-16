dev:
	air

swagger:
	swag init --dir ./,./handlers

testing:
	go test -v ./...