docker-build:
	VERSION=$$(cat VERSION); \
	docker build --build-arg VERSION=$$VERSION -t pluto-server:latest .; \
	docker tag pluto-server:latest registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:latest; \
	docker tag pluto-server:latest registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:$$VERSION

docker-push:
	VERSION=$$(cat VERSION); \
	docker push registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:latest; \
	docker push registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:$$VERSION

docker-run: local-docker-build
	docker run -d -t pluto-server:latest

docker-clean:
	VERSION=$$(cat VERSION); \
	docker rmi registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:latest || true; \
	docker rmi registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:$$VERSION || true; \
	docker rmi pluto-server:latest || true; \
	docker rm -v $(docker ps --filter status=exited -q 2>/dev/null) 2>/dev/null || true; \
	docker rmi $(docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true

binary-build:
	mkdir -p bin; \
	export GO111MODULE=on; \
	go build -o bin/pluto-server cmd/pluto-server/main.go

update-tag:
	VERSION=$$(cat VERSION); \
	git pull --tags; \
	if git --no-pager tag --list | grep $$VERSION -q ; then echo "$$VERSION already exsits"; exit 1; fi; \
	git tag $$VERSION; \
	git push origin $$VERSION

update-changelog:
	VERSION=$$(cat VERSION); \
	git-chglog $$VERSION.. | cat - CHANGELOG.md > temp && mv temp CHANGELOG.md; \
	git commit -am "update CHANGELOG.md"

unit-test:
	export GO111MODULE=on; \
	go test -v ./...

integration-test:
	export GO111MODULE=on; \
	docker-compose -f integration/docker/docker-compose.yml down --rmi all || {return 0}; \
	docker-compose -f integration/docker/docker-compose.yml up --build -d; \
	go run integration/main.go; \
	docker-compose -f integration/docker/docker-compose.yml down --rmi all; \
	docker rmi $(docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true

test: unit-test integration-test

run-server:
	go run cmd/pluto-server/main.go

jenkins-ci: docker-build docker-push docker-clean update-tag update-changelog