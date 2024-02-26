FROM golang:1.22-alpine as golang-builder

RUN addgroup -S components && adduser -S components -G components

WORKDIR /src/app
COPY --chown=components:components . .

RUN go get
RUN go build

FROM scratch
COPY --from=golang-builder /src/app/components /components
COPY --from=golang-builder /etc/passwd /etc/passwd

USER components
CMD [ "/components" ]
