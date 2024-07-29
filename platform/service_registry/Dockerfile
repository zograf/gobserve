FROM golang:1.22

WORKDIR /app

ARG EXPOSE_PORT
EXPOSE $EXPOSE_PORT

COPY ./service_registry ./

RUN go mod download

RUN go build -o run .

#CMD ["./app/run"]
CMD ["sh", "-c", "./run"]