FROM golang:1.18.3-alpine3.16
WORKDIR /app

COPY go.mod go.sum /
COPY . .
RUN go install
RUN go build -o /app/cmd/silkroad

EXPOSE 8080
CMD ["/app/cmd/silkroad"]