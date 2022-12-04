FROM golang:latest

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

COPY . ./

RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]