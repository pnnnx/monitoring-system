FROM golang:1.22-alpine AS build

WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build -o /agent main.go

FROM alpine:3.18
COPY --from=build /agent /agent

CMD ["/agent"]
