FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o bin

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/bin .

EXPOSE 8080 

CMD [ "./bin" ]