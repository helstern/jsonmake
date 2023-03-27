BINARY_NAME=jsonmake
VERSION=0.0.1
RELEASE_NAME=$(VERSION)
RELEASE_DESCRIPTION="Release $(VERSION)"
GITHUB_USERNAME=$(shell git config --get remote.origin.url | sed -n 's/.*\/\/github\.com\/\([^/]*\)\/.*/\1/p')
GITHUB_REPOSITORY=$(shell git config --get remote.origin.url | sed -n 's/.*\/\/github\.com\/[^/]*\/\([^/]*\)\.git.*/\1/p')

build:
	mkdir -p target
	GOOS=linux GOARCH=amd64 go build -o target/$(BINARY_NAME) main.go

tag:
	git tag -a v$(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)


release: build
	curl \
		--header "Authorization: token $(GITHUB_TOKEN)" \
		--header "Content-Type: application/json" \
		--data '{"tag_name": "$(VERSION)", "target_commitish": "v${VERSION}", "name": "$(RELEASE_NAME)", "body": $(RELEASE_DESCRIPTION), "draft": true, "prerelease": false}' \
		https://api.github.com/repos/$(GITHUB_USERNAME)/$(GITHUB_REPOSITORY)/releases
	curl \
		--header "Authorization: token $(GITHUB_TOKEN)" \
		--header "Content-Type: application/octet-stream" \
		--data-binary "@target/$(BINARY_NAME)" \
		https://uploads.github.com/repos/$(GITHUB_USERNAME)/$(GITHUB_REPOSITORY)/releases/$(RELEASE_ID)/assets?name=$(BINARY_NAME)

.PHONY: build release
