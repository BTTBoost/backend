FROM golang:latest
WORKDIR /hashscan
COPY . .
RUN go get -d -v ./...
RUN go install -v ./cmd/fetch_blocks/main.go
ENTRYPOINT ["main"]