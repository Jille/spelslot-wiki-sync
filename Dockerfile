FROM golang:latest as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /serve ./cmd/serve
 
FROM alpine:latest as run

# Copy the application executable from the build image
COPY --from=build /app /charsync-serve

EXPOSE 8080
USER nobody
CMD ["/charsync-serve"]
