# DO NOT USE THIS DOCKERFILE FOR ANY PRODUCTION DEPLOYMENT
# this name is only to test if ipaas is able to correctly detect the dockerfiles in the repo
# and let the user select them correctly`
FROM golang

COPY . .
RUN go mod download
RUN go build -o /app
# as you can see this dockerfile is much more complex than Dockerfile
CMD ["/app"]
