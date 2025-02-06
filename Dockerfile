FROM golang:1.23-alpine
WORKDIR /.
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD [ "/./main" ]
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /.
COPY . .
RUN go build -o main main.go
# Run stage
FROM alpine
WORKDIR /.
COPY --from=builder /app/main .

