FROM golang:1.11 as builder

WORKDIR /go/src/github.com/aaronflower/dzone-shipping/user-cli

COPY . .

# RUN go get -u github.com/golang/dep/cmd/dep
# RUN dep init && dep ensure
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
# COPY user.json /app/user.json
# COPY --from=builder /go/src/github.com/aaronflower/dzone-shipping/service.vessel .
COPY ./user-cli .

# Run the binary as per usual! This time with a binary build in a separate container,
# with all of the correct dependencies and runtime libraries.
ENTRYPOINT ["./user-cli"]
CMD ["./user-cli"]
