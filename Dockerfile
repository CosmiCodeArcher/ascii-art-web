# STAGE 1: builder

FROM golang:1.22.2 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o aaw_app

# STAGE 2: final

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/aaw_app .

COPY templates/ ./templates/
COPY static/ ./static/
COPY banners/ ./banners/

EXPOSE 8081

CMD ["./aaw_app"]