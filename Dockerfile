FROM golang:1.13-alpine AS build

WORKDIR /src

RUN apk add --no-cache git libcap ca-certificates tzdata

# Install go modules
COPY go.mod go.sum ./
RUN go mod download

# Move source
COPY . .

RUN CGO_ENABLED=0 go build -o /static-host

FROM scratch AS run
EXPOSE 80

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/
COPY --from=build /static-host /static-host

ENTRYPOINT ["/static-host"]
