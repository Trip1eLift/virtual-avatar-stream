terraform

```
terraform init
terraform plan
terraform apply
terraform apply -destroy
```

docker build, tag, and push

```
cd ../server

aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 201843717406.dkr.ecr.us-east-1.amazonaws.com

docker build -t virtual-avatar-stream .

docker tag virtual-avatar-stream:latest 201843717406.dkr.ecr.us-east-1.amazonaws.com/virtual-avatar-stream:latest

docker push 201843717406.dkr.ecr.us-east-1.amazonaws.com/virtual-avatar-stream:latest
```

hit websocket with curl
```
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Host: echo.websocket.org" -H "Origin: http://www.websocket.org" http://echo.websocket.org

curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Host: 127.0.0.1:5001" -H "Origin: http://127.0.0.1:5001" http://127.0.0.1:5001

curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Host: 127.0.0.1:5001" -H "Origin: http://127.0.0.1:5001" http://127.0.0.1:5001 || exit 1
```