FROM golang:1.23-alpine AS build

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

ADD cmd/ ./cmd/
ADD internal ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o locker ./cmd/app/

FROM scratch

WORKDIR /app
COPY --from=build /build/locker .

EXPOSE 8080

CMD ["./locker"]