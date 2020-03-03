FROM golang AS build

WORKDIR /workspace
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /static-host


FROM scratch AS run

USER nobody

COPY --from=build /etc_passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /static-host /static-host
