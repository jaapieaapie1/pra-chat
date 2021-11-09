FROM golang:1.17.1

USER root
WORKDIR /xschedule/builddir

COPY . .

RUN go mod download
RUN go build -o /main pra-chat

WORKDIR /

EXPOSE 8080

CMD ["./main"]