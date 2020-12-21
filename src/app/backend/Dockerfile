#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.15.3-alpine3.12 as builder

# gcc
RUN apk add --no-cache build-base

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags acornsoft-dashboard.backend -o acornsoft-dashboard.backend -v backend.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
# FROM honester/httpbin:latest
FROM scratch

COPY --from=builder /usr/src/app/acornsoft-dashboard.backend /app/

ENV KUBECONFIG  "/app/.kube/config"

ENTRYPOINT [ "/app/acornsoft-dashboard.backend" ]

EXPOSE 3001