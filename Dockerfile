FROM golang as go_builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o billing


FROM alpine
WORKDIR /app
COPY --from=go_builder build/billing .

CMD ["./billing"]