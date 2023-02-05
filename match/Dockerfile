FROM golang:1.19.1-alpine3.16

RUN apk --no-cache add curl
RUN apk add --update nodejs npm
RUN npm install -g nodemon

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /virtual-avatar-stream

EXPOSE 5000

# HEALTHCHECK CMD curl -sf http://localhost:5000/health || exit 1 --health-retries=3 --health-timeout=3s --health-interval=60s --health-start-period=5s

CMD [ "/virtual-avatar-stream" ]