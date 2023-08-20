# FROM alpine:3 as downloader
# RUN apk add --no-cache curl

# FROM downloader as rie-downloader
# RUN curl -JLo aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie
# RUN chmod +x aws-lambda-rie
ARG TARGET_PACKAGE
ARG WIRE_TAGS

FROM golang:latest as build
WORKDIR /src
ENV CGO_ENABLED=0
RUN --mount=type=cache,sharing=locked,target=/root/.cache/go-build \
    --mount=type=cache,sharing=locked,target=/go/pkg \
  go install github.com/go-delve/delve/cmd/dlv@latest && \
  go install golang.org/x/vuln/cmd/govulncheck@latest && \
  go install github.com/google/wire/cmd/wire@latest
COPY . .

FROM build as debug-builder
ARG TARGET_PACKAGE
ARG WIRE_TAGS
WORKDIR /src/${TARGET_PACKAGE}
RUN --mount=type=cache,sharing=locked,target=/root/.cache/go-build \
    --mount=type=cache,sharing=locked,target=/go/pkg \
  go generate ./... && \
  wire gen -tags "${WIRE_TAGS}" && \
  go build -gcflags="all=-N -l" -o "/${TARGET_PACKAGE}.debug" .

FROM build as main-builder
ARG TARGET_PACKAGE
ARG WIRE_TAGS
WORKDIR "/src/${TARGET_PACKAGE}"
RUN --mount=type=cache,sharing=locked,target=/root/.cache/go-build \
    --mount=type=cache,sharing=locked,target=/go/pkg \
  go generate ./... && \
  wire gen -tags "${WIRE_TAGS}" && \
  go build -o "/${TARGET_PACKAGE}.run" .

FROM gcr.io/distroless/base:latest as debug
# Note this currently cannot work on MacOS M1.
# - https://stackoverflow.com/a/66370960/1478636
# - https://bytemeta.vip/repo/go-delve/delve/issues/2910
# - https://github.com/docker/for-mac/issues/5191
ARG TARGET_PACKAGE
WORKDIR /opt/main/
COPY --from=build /go/bin/dlv /opt/bin/dlv
COPY --from=debug-builder "/${TARGET_PACKAGE}.debug" "/opt/main/debug.run"
EXPOSE 40000
ENTRYPOINT [ "/opt/bin/dlv" ]
CMD ["--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/opt/main/debug.run"]

FROM gcr.io/distroless/base:latest as main-vanilla
ARG TARGET_PACKAGE
WORKDIR /opt/main/
COPY --from=main-builder "/${TARGET_PACKAGE}.run" "/opt/main/main.run"
ENTRYPOINT [ "/opt/main/main.run" ]

FROM main-vanilla as main-lambda
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.7.0 /lambda-adapter /opt/extensions/lambda-adapter
