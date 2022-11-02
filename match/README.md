## Local Development
```
go mod download
go run ./main.go

npm i -g nodemon
nodemon --watch '*.go' --signal SIGTERM --exec 'go' run main.go
```

## Websocket Testing
wscat
```
npm install -g wscat
wscat -c ws://127.0.0.1:5000/
```

curl ws
```
curl -i -N http://127.0.0.1:5000/ \
-H "Connection: Upgrade" \
-H "Upgrade: websocket" \
-H "Sec-Websocket-Version: 13" \
-H "Sec-Websocket-Key: mock" \
-H "Health: healthcheck" || exit 1
```

curl ws on cmd
```
curl -i -N http://127.0.0.1:5000/ ^
-H "Connection: Upgrade" ^
-H "Upgrade: websocket" ^
-H "Sec-Websocket-Version: 13" ^
-H "Sec-Websocket-Key: mock" ^
-H "Health: healthcheck" || exit 1
```

healthcheck bash
```
curl -i -N http://127.0.0.1:5000/ -H 'Connection: Upgrade' -H 'Upgrade: websocket' -H 'Sec-Websocket-Version: 13' -H 'Sec-Websocket-Key: mock' -H 'Health: healthcheck' || exit 1
```

healthcheck shallow bash
```
curl -i http://127.0.0.1:5000/health || exit 1
```