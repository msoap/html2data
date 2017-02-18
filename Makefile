test:
	go test -cover -race -v $$(glide novendor)

lint:
	golint $$(glide novendor)
	go vet $$(glide novendor)
	errcheck $$(glide novendor)

update-from-github:
	go get -u github.com/msoap/html2data/cmd/html2data

glide-update:
	glide up --update-vendored --strip-vcs --strip-vendor

gometalinter:
	gometalinter --vendor --cyclo-over=25 --line-length=150 --dupl-threshold=150 --min-occurrences=3 --enable=misspell --deadline=10m $$(glide novendor)

generate-manpage:
	docker run -it --rm -v $$PWD:/app -w /app ruby-ronn sh -c 'cat README.md | grep -v "^\[" > html2data.md; ronn html2data.md; mv ./html2data ./html2data.1; rm ./html2data.html ./html2data.md'

create-debian-amd64-package:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w" -o html2data ./cmd/html2data
	set -e ;\
	TAG_NAME=$$(git tag 2>/dev/null | grep -E '^[0-9]+' | tail -1) ;\
	docker run -it --rm -v $$PWD:/app -w /app -e TAG_NAME=$$TAG_NAME ruby-fpm sh -c 'fpm -s dir -t deb --name html2data -v $$TAG_NAME ./html2data=/usr/bin/ ./html2data.1=/usr/share/man/man1/ LICENSE=/usr/share/doc/html2data/copyright README.md=/usr/share/doc/html2data/'
	rm html2data
