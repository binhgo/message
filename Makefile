build:
	go mod vendor
	go mod tidy
	git add vendor -f
	go vet
	go build -mod=vendor