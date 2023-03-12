CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/gorelay
docker stop gorelay
docker rm gorelay
docker build -t gorelay:v1.0 .
sudo docker run -d --name gorelay --restart always --log-opt max-size=10m --log-opt max-file=3 -p 50000:50000  gorelay:v1.0