FROM golang:1.21 AS build_stage

WORKDIR /app
COPY go.mod go.sum ./
COPY pb/*.go ./pb/
COPY *.go ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o rsync-server .

FROM debian:bookworm-slim

COPY --chmod=0700 install_rsync.sh ./
RUN ./install_rsync.sh && rm ./install_rsync.sh
COPY --from=build_stage /app/rsync-server /rsync-server
EXPOSE 25003
ENTRYPOINT ["/rsync-server"]
