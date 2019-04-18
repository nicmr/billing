FROM golang as go_builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY pkg ./pkg
COPY cmd ./cmd

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o billing


FROM alpine
RUN apk add ca-certificates
# Create non-root user
RUN adduser -D runner
# Add AWS config to user
COPY .aws/config /home/runner/.aws/config
# Set variable enabling loading of the config
ENV AWS_SDK_LOAD_CONFIG=1

WORKDIR /app
#copy binaries from build stage
COPY --from=go_builder build/billing .

USER runner
CMD ["./billing cost --month current --bucket altemista-billing"]