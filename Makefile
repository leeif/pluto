VERSION=$(shell cat VERSION)

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t leeif/pluto:latest .
	docker tag leeif/pluto:latest leeif/pluto:$(VERSION)

docker-push:
	docker push leeif/pluto:latest
	docker push leeif/pluto:$(VERSION)

docker-run: local-docker-build
	docker run -d -t pluto-server:latest

docker-clean:
	docker rmi leeif/pluto:latest || true
	docker rmi leeif/pluto:$(VERSION) || true
	docker rm -v $(shell docker ps --filter status=exited -q 2>/dev/null) 2>/dev/null || true
	docker rmi $(shell docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true

binary-build:
	mkdir -p bin
	GO111MODULE=on go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o bin/pluto-server cmd/pluto-server/main.go

check-version-tag:
	git pull --tags
	if git --no-pager tag --list | grep $(VERSION) -q ; then echo "$(VERSION) already exsits"; exit 1; fi

update-tag:
	git pull --tags
	if git --no-pager tag --list | grep $(VERSION) -q ; then echo "$(VERSION) already exsits"; exit 1; fi
	git tag $(VERSION)
	git push origin $(VERSION)

update-changelog:
	git-chglog $(VERSION).. | cat - CHANGELOG.md > temp && mv temp CHANGELOG.md
	git commit -am "update CHANGELOG.md"

unit-test:
	GO111MODULE=on go test -v ./...

test: unit-test

run-server:
	go run cmd/pluto-server/main.go

jenkins-ci: check-version-tag docker-build docker-push docker-clean update-tag update-changelog