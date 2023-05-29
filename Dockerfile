FROM golang:1.19 as BUILD

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build .

FROM alpine

RUN adduser -D -u 1000 -h /app "app"

USER "app"

WORKDIR /app

COPY --from=BUILD /src/outboundip .

CMD [ "./outboundip" ]
