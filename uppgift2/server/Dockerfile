# syntax=docker/dockerfile:1

FROM golang:1.19.5 AS build-stage

WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/* ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /iotuppgift2

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /iotuppgift2 /iotuppgift2

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/iotuppgift2"]