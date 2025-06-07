#docker build --platform linux/amd64,linux/arm64 -t go-server .
#docker build --platform linux/arm64 -t go-server .
docker buildx build  --no-cache --platform arm64 --output=type=docker -t go-server .

docker save go-server >go-server.tar
scp go-server.tar mikrotik:/

