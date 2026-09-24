# STAGE 1: builder

FROM golang:1.22.2 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o aaw_app

# STAGE 2: final

FROM alpine:3.20

LABEL org.opencontainers.image.title="ascii-art-web" \
      org.opencontainers.image.description="Display rendered ascii-art project to a web browser instead of a terminal." \
      org.opencontainers.image.authors="Hamza Ochiponu Musa" \
      org.opencontainers.image.source="https://github.com/cosmicodearcher/ascii-art-web"

WORKDIR /app

RUN adduser -D not_archer

COPY --from=builder /app/aaw_app .

COPY templates/ ./templates/
COPY static/ ./static/
COPY banners/ ./banners/

USER not_archer

EXPOSE 8081

CMD ["./aaw_app"]
