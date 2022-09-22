FROM golang:1.16-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go get
RUN go build -o /out/answer .

FROM scratch AS bin
WORKDIR /app
COPY --from=build /out/answer /app/
CMD ["/app/answer"]
