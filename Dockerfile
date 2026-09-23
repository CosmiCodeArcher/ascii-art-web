FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN go build -o aaw_app

EXPOSE 8081

CMD ["./aaw_app"]