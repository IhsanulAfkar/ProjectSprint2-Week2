FROM golang:1.22

# RUN apk update && apk add --no-cache git
# RUN pwd
WORKDIR /Week2

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /eniqilo
EXPOSE 8080
CMD ["/eniqilo"]