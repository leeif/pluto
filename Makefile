VERSION=$(shell cat VERSION)

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t mushare/pluto:latest .
	docker tag mushare/pluto:latest mushare/pluto:$(VERSION)

docker-build-staging:
	docker build --build-arg VERSION=staging -t mushare/pluto:staging .

docker-push:
	docker push mushare/pluto:latest
	docker push mushare/pluto:$(VERSION)

docker-push-staging:
	docker push mushare/pluto:staging

docker-clean:
	docker rmi mushare/pluto:latest || true
	docker rmi mushare/pluto:$(VERSION) || true
	docker rm -v $(shell docker ps --filter status=exited -q 2>/dev/null) 2>/dev/null || true
	docker rmi $(shell docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true

docker-clean-staging:
	docker rmi mushare/pluto:staging || true
	docker rm -v $(shell docker ps --filter status=exited -q 2>/dev/null) 2>/dev/null || true
	docker rmi $(shell docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true

check-version-tag:
	git pull --tags
	if git --no-pager tag --list | grep $(VERSION) -q ; then echo "$(VERSION) already exsits"; exit 1; fi

update-tag:
	git pull --tags
	if git --no-pager tag --list | grep $(VERSION) -q ; then echo "$(VERSION) already exsits"; exit 1; fi
	git tag $(VERSION)
	git push origin $(VERSION)

server-binary-build:
	mkdir -p bin
	GO111MODULE=on go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o bin/pluto-server cmd/pluto-server/main.go

migrate-binary-build:
	mkdir -p bin
	GO111MODULE=on go build -o bin/pluto-migrate cmd/pluto-migrate/main.go

unit-test:
	GO111MODULE=on go test -v ./...

test: unit-test

ci-build-production: test check-version-tag docker-build docker-push docker-clean update-tag

ci-build-staging: test docker-build-staging docker-push-staging docker-clean-staging
