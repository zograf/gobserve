FROM golang:1.22

WORKDIR /app

COPY ./proxy ./

RUN go mod download

RUN go build -o run .

#CMD ["./app/run"]
CMD ["sh", "-c", "./run"]