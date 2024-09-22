FROM --platform=linux/amd64 golang:1.21.3

RUN apt-get update && apt-get install -y default-mysql-client

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

RUN go mod download

RUN go install github.com/cosmtrek/air@v1.29.0

WORKDIR /app

COPY . .

COPY entrypoint.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]

WORKDIR /app
CMD ["air", "-c", ".air.toml"]