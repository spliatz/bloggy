# Build
FROM golang:alpine AS build

WORKDIR /go/src/bloggy-backend/
COPY . .
RUN go mod download
RUN go build -o bloggy_backend cmd/app.go

# Deploy
FROM alpine:latest
COPY --from=build /go/src/bloggy-backend/bloggy_backend .

ENTRYPOINT [ "./bloggy_backend" ]
