FROM golang:1.22.2 as build-env
# All these steps will be cached
RUN mkdir /ApiServer
WORKDIR /ApiServer
COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/ApiServer
FROM scratch 
COPY --from=build-env /go/bin/ApiServer /go/bin/ApiServer
ENTRYPOINT ["/go/bin/ApiServer"]