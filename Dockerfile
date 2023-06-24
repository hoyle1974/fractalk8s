FROM golang:1.20-bullseye as base

WORKDIR /app

COPY ./service/go.mod /app/service/go.mod
COPY ./service/go.sum /app/service/go.sum
WORKDIR /app/service
RUN go mod download

COPY . /app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o /cmd

FROM gcr.io/distroless/static-debian11

COPY --from=base /cmd .
COPY --from=base /app/ /app/


CMD [ "/cmd" ]
