FROM golang:1.11
WORKDIR /service
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o service_monitor main.go
# ENTRYPOINT [ "./service_monitor" ]


# docker build -t esc .
#  docker run -p 5000:5000 --name esc -t esc 