#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.15.3-alpine3.12 as builder

# gcc
RUN apk add --no-cache build-base

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags acornsoftlab.kore3.dashboard -o acornsoftlab.kore3.dashboard -v src/app/backend/dashboard.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
# FROM honester/httpbin:latest
FROM alpine

COPY --from=builder /usr/src/app/acornsoftlab.kore3.dashboard /app/
COPY --from=builder /usr/src/app/i18n/                        /app/public/
COPY --from=builder /usr/src/app/i18n/locale_conf.json        /app/

WORKDIR /app

ENV KUBECONFIG /app/.kube/config
ENV INSECURE_BIND_ADDR 0.0.0.0
ENV BIND_ADDR 0.0.0.0
ENV SIDECAR_HOST http://localhost:8000

ENTRYPOINT /app/acornsoftlab.kore3.dashboard --kubeconfig=${KUBECONFIG} --insecure-bind-address=${INSECURE_BIND_ADDR} --bind-address=${BIND_ADDR} --sidecar-host=${SIDECAR_HOST}

EXPOSE 9090