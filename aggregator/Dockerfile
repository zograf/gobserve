FROM golang:1.22

WORKDIR /app

ARG EXPOSE_PORT
EXPOSE $EXPOSE_PORT

ARG PROJECT_DIR
COPY ./$PROJECT_DIR ./

RUN go mod download

RUN go build -o run .

#CMD ["./app/run"]
CMD ["sh", "-c", "./run"]