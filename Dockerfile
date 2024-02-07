FROM golang:latest
ENV GOPATH=/
COPY . .
RUN go get -d -v ./...
RUN go build -o wallet-infotecs ./cmd/app/main.go
EXPOSE 8000
CMD ["./wallet-infotecs"]