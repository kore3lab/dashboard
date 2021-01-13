#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.15.3-alpine3.12 as builder

# gcc
RUN apk add --no-cache build-base

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags acornsoft-dashboard.dashboard -o acornsoft-dashboard.dashboard -v src/app/backend/dashboard.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
# FROM honester/httpbin:latest
FROM alpine

COPY --from=builder /usr/src/app/acornsoft-dashboard.dashboard	/app/
COPY --from=builder /usr/src/app/i18n/                          /app/public/
COPY --from=builder /usr/src/app/i18n/locale_conf.json          /app/

WORKDIR /app

ENV KUBECONFIG ""

ENTRYPOINT ["/app/acornsoft-dashboard.dashboard", "--insecure-bind-address=0.0.0.0", "--bind-address=0.0.0.0"]
EXPOSE 9090