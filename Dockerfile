FROM golang:1.22.3

WORKDIR /var/www
COPY . /var/www

RUN go mod download && go mod verify