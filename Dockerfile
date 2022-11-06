FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o test_start cmd/main.go

ENTRYPOINT ["./test_start"]

#docker build . -t sber-test-task
#docker run -t sber-test-task -v -u 185.204.3.165