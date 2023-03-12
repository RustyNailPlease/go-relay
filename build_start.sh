GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo  -o ./out/gorelay
chmod +x out/gorelay
sudo docker stop gorelay
sudo docker rm gorelay
sudo docker build --no-cache -t gorelay:v1.0 .
sudo docker run -d --name gorelay --restart always --log-opt max-size=10m --log-opt max-file=3 -p 50000:50000  gorelay:v1.0