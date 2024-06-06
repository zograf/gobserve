FROM golang:1.22

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./
#COPY core ./
#COPY proxy ./
#COPY service_registry ./
#COPY simple_proxy ./
#COPY main.go ./
COPY . ./

RUN go mod download

WORKDIR /app/proxy

RUN go build -o run .

CMD ["run"]