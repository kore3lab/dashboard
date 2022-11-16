#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.19-alpine as builder

ADD . /usr/src/app
WORKDIR /usr/src/app

# gcc
RUN apk add --no-cache build-base

# build
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags backend -o backend -v backend.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM ubuntu:18.04

COPY --from=builder /usr/src/app/backend /app/

RUN apt-get update && apt-get install -y wget
RUN wget https://t1.kakaocdn.net/kakaoicloud/binaries-kic-iam-auth/v0.1.2/kic-iam-auth_0.1.2_Linux_x86_64/kic-iam-auth && chmod +x kic-iam-auth && mv kic-iam-auth /usr/local/bin/

#-------------------------------------------
# STEP 3 : execute a binary
#-------------------------------------------
ARG app_env=production
ENV KUBECONFIG  ""
ENV GIN_MODE=release
ENV APP_ENV=${app_env:-"production"}

ENTRYPOINT [ "/app/backend" ]

EXPOSE 3001
