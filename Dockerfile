FROM golang:1.22 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM gcr.io/distroless/static:nonroot

ENV PORT=8080
EXPOSE 8080

WORKDIR /
COPY --from=build /out/api /api

USER nonroot:nonroot
ENTRYPOINT ["/api"]

