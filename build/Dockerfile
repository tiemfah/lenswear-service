FROM golang:1.17-alpine AS build_base
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build ./cmd/api/main.go

FROM alpine:3.15
COPY --from=build_base /app/main .
COPY --from=build_base /app/configs/local-config.yaml ./configs/local-config.yaml
COPY --from=build_base /app/cert ./cert
CMD ./main -env=local