FROM golang:1.20.3-alpine3.16 as build

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o /kanbango-gateway ./cmd/main.go

FROM alpine:3.16
COPY --from=build /kanbango-gateway /kanbango-gateway
COPY ./kbg-cert.pem /
COPY ./kbg-key.pem /
# COPY --from=build /usr/src/app/.env /.env

CMD ["/kanbango-gateway"]