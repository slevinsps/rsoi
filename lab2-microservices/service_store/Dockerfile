FROM golang:1.11
WORKDIR /service
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o service_store main.go
ENTRYPOINT [ "./service_store" ]

# docker build -t esc .
#  docker run -p 5000:5000 --name esc -t esc 