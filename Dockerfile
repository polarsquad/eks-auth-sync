FROM golang:1.14 AS build
WORKDIR /build
COPY go.mod ./
COPY bin/ bin/
COPY cmd/ cmd/
COPY internal/ internal/
ARG APP_VERSION
ARG GIT_HASH
RUN ./bin/build.sh

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /build/eks-auth-sync /usr/local/bin/eks-auth-sync
RUN eks-auth-sync --version
ENTRYPOINT ["/usr/local/bin/eks-auth-sync"]
