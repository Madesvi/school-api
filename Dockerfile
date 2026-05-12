# --- Stage 1: Builder ---
FROM golang:1.25.7-alpine AS builder

RUN apk add --no-cache curl make nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN templ generate
RUN npx @tailwindcss/cli -i ./views/css/input.css -o ./public/output.css --minify
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o main .

# --- Stage 2: Runner ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
ENV APP_VERSION=latest
EXPOSE 3000
CMD ["./main"]
