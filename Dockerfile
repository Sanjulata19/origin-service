FROM golang:1.13-alpine AS build

WORKDIR /src

RUN apk add --no-cache git libcap ca-certificates
RUN update-ca-certificates 2>/dev/null || true

# Install go modules
COPY go.mod go.sum ./
RUN go mod download

# Move source
COPY . .

RUN CGO_ENABLED=0 go build -o /static-host

FROM scratch AS run
EXPOSE 80

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /static-host /static-host

ENTRYPOINT ["/static-host"]
