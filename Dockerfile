FROM golang:1.20-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go mod tidy
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN task build

FROM scratch

COPY --from=builder ["/build/http-server", "/http-server"]

ENV GO_ENV=production

CMD ["/http-server"]

