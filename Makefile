APP_NAME := html2data
APP_DESCRIPTION := $$(awk 'NR == 11' README.md)
APP_URL := https://github.com/msoap/$(APP_NAME)
APP_MAINTAINER := $$(git show HEAD | awk '$$1 == "Author:" {print $$2 " " $$3 " " $$4}')
GIT_TAG := $$(git describe --tags --abbrev=0)

test:
	go test -cover -race -v ./...

lint:
	golint $$(glide novendor)
	go vet ./...
	errcheck $$(glide novendor)

update-from-github:
	go get -u github.com/msoap/$(APP_NAME)/cmd/$(APP_NAME)

dep-ensure:
	dep ensure

gometalinter:
	gometalinter --vendor --cyclo-over=25 --line-length=150 --dupl-threshold=150 --min-occurrences=3 --enable=misspell --deadline=10m --exclude=SA1022 $$(glide novendor)

generate-manpage:
	cat README.md | grep -v "^\[" > $(APP_NAME).md
	docker run --rm -v $$PWD:/app -w /app msoap/ruby-ronn ronn $(APP_NAME).md
	mv ./$(APP_NAME) ./$(APP_NAME).1
	rm ./$(APP_NAME).{md,html}

create-debian-amd64-package:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(APP_NAME) ./cmd/html2data
	docker run --rm -v $$PWD:/app -w /app msoap/ruby-fpm \
		fpm -s dir -t deb --force --name $(APP_NAME) -v $(GIT_TAG) \
			--license="$$(head -1 LICENSE)" \
			--url=$(APP_URL) \
			--description="$(APP_DESCRIPTION)" \
			--maintainer="$(APP_MAINTAINER)" \
			--category=network \
			./$(APP_NAME)=/usr/bin/ \
			./$(APP_NAME).1=/usr/share/man/man1/ \
			LICENSE=/usr/share/doc/$(APP_NAME)/copyright \
			README.md=/usr/share/doc/$(APP_NAME)/
	rm $(APP_NAME)
