# Stage 1. Build app
ARG GO_VERSION=1.18

FROM golang:${GO_VERSION} AS build

WORKDIR /usr/local/go/src/app/

COPY ./cmd/ /usr/local/go/src/app/cmd/
COPY ./internal/ /usr/local/go/src/app/internal/
COPY ./pkg/ /usr/local/go/src/app/pkg/
COPY config.yml /usr/local/go/src/app/
COPY go.mod /usr/local/go/src/app/
COPY go.sum /usr/local/go/src/app/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app /usr/local/go/src/app/cmd/main/app.go

# STAGE 2: build the container to run
FROM scratch AS final

WORKDIR /root/

# copy compiled app
COPY --from=build /app ./
COPY --from=build /usr/local/go/src/app/config.yml ./

CMD ["./app"]