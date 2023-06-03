FROM golang:latest

# RUN mkdir /build
ADD . /build
WORKDIR /build

# COPY go.mod go.sum ./
# COPY * ./
RUN go mod download
RUN go install

# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go-server
RUN go build -o /docker-go-server

EXPOSE 8080

CMD ["/docker-go-server"]