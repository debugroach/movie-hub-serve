# Build stage
FROM golang:1.22-alpine3.19 
WORKDIR /app
COPY . .
RUN go build -o main main.go
CMD [ "/app/main" ]
EXPOSE 8080 8080
