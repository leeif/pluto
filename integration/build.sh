docker-compose -f docker/docker-compose.yml up --build -d

export GO111MODULE=on

go run .

docker rmi $(docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true