## Pre-requisites

1. terraform CLI
2. aws CLI
3. docker CLI
4. psql CLI


## Terraform command

```shell
terraform init
terraform plan
terraform apply
terraform apply -destroy
```

docker build, tag, and push is now part of terraform.

## Stack basic testing

```shell
curl http://virtualavatar-stream.trip1elift.com:80/health -L
```

```shell
curl https://virtualavatar-stream.trip1elift.com:443/health
```

```shell
curl http://virtualavatar-stream.trip1elift.com/health -L
```

```shell
curl https://virtualavatar-stream.trip1elift.com/health
```

## Stack database testing

```shell
curl https://virtualavatar-stream.trip1elift.com:443/health-database
```

```shell
curl https://virtualavatar-stream.trip1elift.com:443/health-proxy
```