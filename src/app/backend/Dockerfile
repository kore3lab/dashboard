#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.19-alpine as builder

# gcc
RUN apk add --no-cache build-base

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags backend -o backend -v backend.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM scratch

COPY --from=builder /usr/src/app/backend /app/

#-------------------------------------------
# STEP 3 : execute a binary
#-------------------------------------------
ARG app_env=production
ENV KUBECONFIG  ""
ENV GIN_MODE=release
ENV APP_ENV=${app_env:-"production"}

ENTRYPOINT [ "/app/backend" ]

EXPOSE 3001
