FROM golang:1.15-buster as compile-image

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch

COPY --from=compile-image /app/app /

CMD ["/app"]
