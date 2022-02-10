FROM golang:latest
WORKDIR /hashscan
COPY . .
RUN go get -d -v ./...
RUN go install -v ./cmd/fetch_opensea_trades/main.go
ENTRYPOINT ["main"]