FROM golang:alpine AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/server/main.go

FROM alpine:latest
COPY --from=build /app /app
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app"]
