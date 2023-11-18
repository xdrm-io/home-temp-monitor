FROM golang:buster as builder

ARG arg_os=linux
ARG arg_arch=amd64

# copy sources
RUN mkdir -p /app
WORKDIR /app

# for caching 'go mod download'
COPY go.sum .
COPY go.mod .
RUN go mod download

# copy sources
COPY . .

# compile
RUN CGO_ENABLED=0 GOOS=${arg_os} GOARCH=${arg_arch} go build -ldflags="-s -w" -o /app/binary .

FROM gcr.io/distroless/static-debian11 as production
COPY --from=builder /app/binary /app/
WORKDIR /app
EXPOSE 80/tcp
ENTRYPOINT [ "/app/binary" ]