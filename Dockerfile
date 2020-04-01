FROM golang:alpine
RUN apk add --no-cache git

WORKDIR /go/src/github.com/otz1/scraper
COPY . .

RUN go get ./...
RUN go install

CMD ["scraper"]

EXPOSE 8001