FROM golang

COPY . .
RUN go mod download
RU go build -o /app

CMD ["/app"]
