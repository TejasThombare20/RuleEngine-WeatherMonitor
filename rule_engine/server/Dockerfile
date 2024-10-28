FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

EXPOSE 8000
CMD ["./main"]

# Note: Even though we EXPOSE 8000, this service will be accessible 
# on port 8000 on the host machine due to the port mapping in docker-compose.yml