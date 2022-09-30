## Python server local run
```
cd server
python3 -m venv ./venv
```

Activate virtual environment
```
./venv/Scripts/activate     (Windows)
source ./venv/bin/activate  (Mac)
```

operations
```
pip list
pip install -r requirements.txt
python main.py
deactivate
```

## hit websocket with curl
```
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Host: echo.websocket.org" -H "Origin: http://www.websocket.org" http://echo.websocket.org || exit 1

curl -i -N \
-H "Connection: close" \
-H "Upgrade: websocket" \
-H "Host: 127.0.0.1:5001" \
-H "Origin: http://127.0.0.1:5001" \
-H "Sec-WebSocket-Version: 13" \
-H "Sec-Websocket-Key: healthcheck" \
http://127.0.0.1:5001

curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Host: 127.0.0.1:5001" -H "Origin: http://127.0.0.1:5001" http://127.0.0.1:5001 || exit 1
```