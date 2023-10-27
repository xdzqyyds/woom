FROM node:20-alpine as builder-node

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN npm run build

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY --from=builder-node /app/dist /app/dist

RUN go build -o woom

# Bin
FROM alpine AS bin

COPY --from=builder /app/woom /usr/bin/woom

EXPOSE 8080/tcp

CMD ["-migrate"]

ENTRYPOINT ["/usr/bin/woom"]