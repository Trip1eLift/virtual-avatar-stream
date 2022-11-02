FROM node:18-alpine3.15

WORKDIR /app

COPY package.json .

RUN npm install

COPY . .

EXPOSE 3000