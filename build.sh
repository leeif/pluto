set -e

version=`cat VERSION`

echo "version $version"

docker build --build-arg VERSION=$version -t pluto-server:latest .
docker tag pluto-server:latest registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:latest
docker tag pluto-server:latest registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:${version}


docker push registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:latest
docker push registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:${version}

echo "clean docker"
docker rmi registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:latest
docker rmi registry.cn-hongkong.aliyuncs.com/mushare/pluto-server:${version}
docker rmi pluto-server:latest
docker rm -v $(docker ps --filter status=exited -q 2>/dev/null) 2>/dev/null || true
docker rmi $(docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true