FROM golang:1.20-alpine

RUN addgroup -S components && adduser -S components -G components
USER components

WORKDIR /src/app
COPY --chown=components:components . .

RUN go get
RUN go install
