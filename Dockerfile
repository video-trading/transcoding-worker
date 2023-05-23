# Build stage
FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Production stage
FROM linuxserver/ffmpeg

WORKDIR /root/

COPY --from=build /app/app .

ENTRYPOINT ["./app"]
CMD ["./app"]