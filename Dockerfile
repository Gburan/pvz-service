FROM golang:1.24.2-alpine AS build_stage
WORKDIR /go/bin/app_build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app_build ./cmd/main.go


FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /app_build ./app_build
COPY config/config.yaml ./config/config.yaml
COPY migrations ./migrations
RUN chmod +x ./app_build
EXPOSE 8080/tcp
ENTRYPOINT ["./app_build"]
