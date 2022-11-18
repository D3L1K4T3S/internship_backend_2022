FROM golang:1.19 AS builder

WORKDIR /cmd/main

COPY .. .

RUN go mod download -x && go mod verify
RUN cd cmd/main && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -v -o /cmd/main/build main.go

FROM scratch
COPY --from=builder /cmd/main/build /cmd/main/build