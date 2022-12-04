FROM golang:latest
RUN apt update
RUN apt install -y dnsutils
WORKDIR /app
COPY ./app src
RUN cd src && if [ ! -d "vendor" ]; then echo "go dep not found, will download" && go mod download -x; fi
RUN cd src && go build -o ../server server/main.go
CMD ["./server"]