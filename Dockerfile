FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /bin/jenkins-build-info ./cmd/jenkins-build-info

FROM alpine:3.18
COPY --from=builder /bin/jenkins-build-info /bin/jenkins-build-info
ENTRYPOINT ["/bin/jenkins-build-info"]