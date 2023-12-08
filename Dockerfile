ARG base_image=golang:1.21

FROM ${base_image} as build

ENV GO111MODULE=on

RUN apt update && apt install ca-certificates libgnutls30 -y

RUN mkdir -p /go/src

WORKDIR /go/src

COPY . .

RUN set -xe \
    && go mod tidy \
    && go mod vendor -v \
    && go build -ldflags "-linkmode external -extldflags -static" -a main.go

FROM scratch
COPY --from=build /go/src/main /main
ENV PORT=8080
ENTRYPOINT ["/main"]

LABEL image.name="webfinger" \
      image.description="WebFinger server"
