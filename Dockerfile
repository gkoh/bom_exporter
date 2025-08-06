FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum /app/
COPY bom /app/bom/
COPY cmd /app/cmd/
RUN go build cmd/bom_exporter.go

FROM alpine
WORKDIR /
COPY --from=build /app/bom_exporter /bom_exporter
EXPOSE 8080
ENTRYPOINT ["/bom_exporter"]
