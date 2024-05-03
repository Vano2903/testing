FROM golang

COPY . .
RUN go mod download
RUN go build -o /app
CMD ["/app"]
