# Choose whatever you want, version >= 1.16
FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ENV GOOS=linux
ENV GOARCH=amd64

CMD ["go", "build", "-o", "murkoto/eniqilo", "."]
