FROM golang as go_builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o billing


FROM alpine
RUN apk add ca-certificates
# Create non-root user
RUN adduser -D runner
# Add AWS config to user
COPY .aws /home/runner/.aws
# Set variable enabling loading of the config
ENV AWS_SDK_LOAD_CONFIG=1

WORKDIR /app
#copy binaries from build stage
COPY --from=go_builder build/billing .

USER runner
CMD ["./billing"]