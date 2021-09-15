FROM golang:1.11
WORKDIR /db-forum
COPY ./restService/go.mod . 
COPY ./restService/go.sum .
RUN go mod download
COPY . .
RUN cd ./restService && go build -o main .
ENTRYPOINT [ "./restService/main" ]

# docker build -t esc .
#  docker run -p 5000:5000 --name esc -t esc 