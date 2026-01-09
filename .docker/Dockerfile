FROM golang:1.25.5

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@v1.52.3

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./cmd/main.go

CMD ["air", "-c", ".air.toml"]