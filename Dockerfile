FROM golang:1.25 AS base
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o llm-inference
EXPOSE 5080
CMD ["/build/llm-inference"]