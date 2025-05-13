FROM golang:1.24.3-nanoserver-ltsc2025
WORKDIR /usr/local/app

RUN apt-get update &&\
    apt-get upgrade &&\
    apt-get install -y gcc libx11-dev libgl1-mesa-dev xorg-dev xvfb

COPY . .

RUN go mod tidy
RUN go build -v main.go

CMD ["./main"]