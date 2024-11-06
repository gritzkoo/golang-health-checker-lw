test:
	go test -race -coverprofile=profile.cov ./...
coverage: test
	go tool cover -html=profile.cov
build:
	go build -race -a -installsuffix cgo -o healthchecker pkg/**/*.go
view-docs:
	godoc -http=:8331
run:
	go run main.go
install-godoc:
	go install golang.org/x/tools/cmd/godoc@latest
