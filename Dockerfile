FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /nodeapp

EXPOSE 16574

CMD [ "/nodeapp", "s" ]