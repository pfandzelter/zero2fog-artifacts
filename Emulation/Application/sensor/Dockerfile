FROM golang:latest AS build

WORKDIR /go/src/sensor

COPY . .

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go install ./...

# actual Docker image
FROM scratch

WORKDIR /

COPY --from=build /go/bin/sensor bin

ENTRYPOINT ["./bin"]
