FROM golang:latest

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -o web .

CMD ["/app/web"]
