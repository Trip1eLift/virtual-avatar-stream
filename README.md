# Virtual-Avatar-Streaming-Backend
This is a PoC project to test how to deploy a scalable end-to-end streaming distributive system to support video chat.

## Python server docker run
```
cd server
docker build --tag websocket-server .; docker run -it -p 5001:5001 websocket-server (Powershell)
```

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

## Client
```
cd client
npm install
npm start
```

## Fargate deployment docs
1. https://section411.com/2019/07/hello-world/
2. https://engineering.finleap.com/posts/2020-02-20-ecs-fargate-terraform/