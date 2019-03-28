FROM golang as go_builder
WORKDIR /build

COPY main.go .
RUN go get -d "github.com/aws/aws-sdk-go/aws"
RUN CGO_ENABLED=0 GOOS=linux go build -a -o billing

FROM alpine
WORKDIR /app
COPY --from=go_builder /build/billing .
CMD ["./billing"]