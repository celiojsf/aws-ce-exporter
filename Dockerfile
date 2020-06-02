FROM golang:1.14 as builder
WORKDIR /go/src/github.com/celiojsf/aws-ce-exporter
COPY . .
#RUN make setup
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/aws-ce-exporter

FROM ubuntu
EXPOSE 2112
ENV TZ=America/Sao_Paulo
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /tmp/aws-ce-exporter .
ENTRYPOINT ["./aws-ce-exporter"]