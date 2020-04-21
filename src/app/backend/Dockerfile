#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.15.3-alpine3.12 as builder

# gcc
RUN apk add --no-cache build-base

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags acornsoftlab.kore3.backend -o acornsoftlab.kore3.backend -v backend.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
# FROM honester/httpbin:latest
FROM scratch

COPY --from=builder /usr/src/app/acornsoftlab.kore3.backend /app/

ENV KUBECONFIG  "/app/.kube/config"

ENTRYPOINT [ "/app/acornsoftlab.kore3.backend" ]

EXPOSE 3001