FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o expense-tracker-api .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& mkdir -p /app/data

COPY --from=builder /app/expense-tracker-api ./expense-tracker-api
COPY --from=builder /app/conf ./conf
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/data/.gitkeep ./data/.gitkeep

EXPOSE 8080

CMD ["./expense-tracker-api"]
