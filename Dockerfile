FROM golang:1.8.0
EXPOSE 8080
CMD gennaker start
COPY . /go/src/github.com/vgheri/gennaker
WORKDIR /go/src/github.com/vgheri/gennaker
RUN make install
