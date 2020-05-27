FROM golang:latest AS build

WORKDIR /go/src/cfd

COPY . .

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go install ./...


# actual Docker image
FROM scratch

WORKDIR /

COPY --from=build /go/bin/cfd bin

EXPOSE 8080

ENTRYPOINT ["./bin"]
