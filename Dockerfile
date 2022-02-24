FROM golang:1.14-alpine

RUN apk add --no-cache git

WORKDIR /app/news-app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/news-app .

EXPOSE 8000

CMD ["./out/news-app"]