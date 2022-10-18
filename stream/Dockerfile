FROM golang:1.19.1-alpine3.16

RUN apk add --update nodejs npm

RUN npm install -g nodemon

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /virtual-avatar-stream

EXPOSE 5001

CMD ["/virtual-avatar-stream"]