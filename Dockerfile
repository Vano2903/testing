FROM go:1.22

COPY . .
RUN go mod download
RUN go build -o /app

CMD ["/app"]
