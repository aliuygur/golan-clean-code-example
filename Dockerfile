FROM golang:1.16-alpine

# set default timezone
ENV TZ Europe/Istanbul
RUN apk update && apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go install ./cmd/...

CMD ["restserver"]