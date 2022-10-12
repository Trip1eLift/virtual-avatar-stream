## Local Development
```
go mod download
go run ./main.go

npm i -g nodemon
nodemon --watch '*.go' --signal SIGTERM --exec 'go' run main.go
```

launch go from docker
```
docker build --tag virtual-avatar-stream .; docker run -it -p 5001:5001 virtual-avatar-stream
```

launch with health check (healthcheck defined in Dockerfile)
```
docker build --tag virtual-avatar-stream .; docker run -d -p 5001:5001 virtual-avatar-stream

docker inspect --format='{{json .State.Health}}' 221936098d1362afb92ee0c3397c4d2aa2f8cb9d59807c627a06b0f4b558d153
```

## Push to AWS ECR
```
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 201843717406.dkr.ecr.us-east-1.amazonaws.com

docker build -t virtual-avatar-stream .

docker tag virtual-avatar-stream:latest 201843717406.dkr.ecr.us-east-1.amazonaws.com/virtual-avatar-stream:latest

docker push 201843717406.dkr.ecr.us-east-1.amazonaws.com/virtual-avatar-stream:latest
```

## Websocket Testing
wscat
```
npm install -g wscat
wscat -c ws://127.0.0.1:5001/
```

curl ws
```
curl -i -N http://127.0.0.1:5001/ \
-H "Connection: Upgrade" \
-H "Upgrade: websocket" \
-H "Sec-Websocket-Version: 13" \
-H "Sec-Websocket-Key: mock" \
-H "Health: healthcheck" || exit 1
```

curl ws on cmd
```
curl -i -N http://127.0.0.1:5001/ ^
-H "Connection: Upgrade" ^
-H "Upgrade: websocket" ^
-H "Sec-Websocket-Version: 13" ^
-H "Sec-Websocket-Key: mock" ^
-H "Health: healthcheck" || exit 1
```

healthcheck bash
```
curl -i -N http://127.0.0.1:5001/ -H 'Connection: Upgrade' -H 'Upgrade: websocket' -H 'Sec-Websocket-Version: 13' -H 'Sec-Websocket-Key: mock' -H 'Health: healthcheck' || exit 1
```

healthcheck shallow bash
```
curl -i http://127.0.0.1:5001/health || exit 1
```