FROM golang:1.14 AS build
WORKDIR /build
COPY go.mod ./
COPY cmd/ cmd/
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o eks-auth-sync ./cmd/eksauthsync 

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /build/eks-auth-sync /usr/local/bin/eks-auth-sync
ENTRYPOINT ["/usr/local/bin/eks-auth-sync"]
