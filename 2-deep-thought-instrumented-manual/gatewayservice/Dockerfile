FROM golang:1.16-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go get
RUN go build -o /out/gatewayservice .

FROM scratch AS bin
WORKDIR /app
COPY --from=build /out/gatewayservice /app/
CMD ["/app/gatewayservice"]
