test:
	go test -coverprofile=profile.cov ./...
coverage: test
	go tool cover -html=profile.cov
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o healthchecker pkg/**/*.go
view-docs:
	godoc -http=:8331
run:
	go run main.go
