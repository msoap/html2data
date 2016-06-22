test:
	go test -cover -race -v $(glide novendor)

lint:
	golint $(glide novendor)
	go vet $(glide novendor)
	errcheck $(glide novendor)

update-from-github:
	go get -u github.com/msoap/html2data/cmd/html2data
