FROM golang:1.13.8

WORKDIR /go/src/github.com/ZdravkoGyurov/go-web-app
COPY . .

RUN go get -d -v go.mongodb.org/mongo-driver/mongo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

EXPOSE 8080

CMD ["./main"]