FROM golang:1.23-alpine as golang-builder

RUN addgroup -S components && adduser -S components -G components

WORKDIR /src/app
COPY --chown=components:components . .

RUN go get
RUN go build

FROM scratch
COPY --from=golang-builder /src/app/components /components
COPY --from=golang-builder /etc/passwd /etc/passwd
COPY --chown=components:components static /static

USER components
CMD [ "/components" ]
