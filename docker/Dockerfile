FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

EXPOSE 5000

RUN go build -o expense ./cmd/api

CMD ["./expense"]