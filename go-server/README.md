launch go from local
```
go mod download
go run ./main.go
```

launch go from docker
```
docker build --tag virtual-avatar-stream .; docker run -it -p 5001:5001 virtual-avatar-stream
```

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