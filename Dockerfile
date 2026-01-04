FROM golang:alpine AS build-stage

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /lendbook ./cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker/main.go

FROM gcr.io/distroless/static-debian13 AS build-release

WORKDIR /

COPY --from=build-stage /lendbook /lendbook

COPY --from=build-stage /worker /worker

COPY db/migrations ./db/migrations

USER nonroot:nonroot

ENTRYPOINT ["/lendbook"]