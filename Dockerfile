FROM golang:1.19-alpine

RUN addgroup -S components && adduser -S components -G components
USER components

WORKDIR /src/app
COPY . .

RUN go get
RUN go install
