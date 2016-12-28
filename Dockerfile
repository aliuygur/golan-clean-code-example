FROM golang:latest
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app
RUN GOPATH=`pwd` go get -v app/...
EXPOSE 5000 
ENTRYPOINT ["/app/bin/api"]