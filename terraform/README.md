terraform

```
terraform init
terraform plan
terraform apply
terraform apply -destroy
```

docker build, tag, and push is now part of terraform.

## Trouble shooting

1. CannotPullContainerError: inspect image has been retried 1 time(s): failed to resolve ref "docker.io/library/virtual-avatar-stream:latest": pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed
https://aws.amazon.com/premiumsupport/knowledge-center/ecs-pull-container-api-error-ecr/

2. CannotPullContainerError: inspect image has been retried 1 time(s): failed to resolve ref "201843717406.dkr.ecr.us-east-1.amazonaws.com/virtual-avatar-stream:latest": 201843717406.dkr.ecr.us-east-1.amazonaws.com/virtual-avatar-stream:latest: not found