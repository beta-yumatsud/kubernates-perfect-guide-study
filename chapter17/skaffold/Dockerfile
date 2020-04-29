# Stage 1
FROM golang:1.13-alpine3.10 as builder
COPY ./main.go ./
RUN go build -o /go-app ./main.go

# Stage 2
FROM alpine:latest
EXPOSE 8080
COPY --from=builder /go-app .
ENTRYPOINT ["./go-app"]
