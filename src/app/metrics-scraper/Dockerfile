#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.19-alpine as builder

# gcc
RUN apk add --no-cache build-base sqlite-libs sqlite-dev sqlite

ADD . /usr/src/app
WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=1 go build -ldflags '-w -extldflags "-static"' -tags dashboard-metrics-scraper -o dashboard-metrics-scraper -v server.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM ubuntu:18.04

COPY --from=builder /usr/src/app/dashboard-metrics-scraper /app/

RUN apt-get update && apt-get install -y wget
RUN wget https://t1.kakaocdn.net/kakaoicloud/binaries-kic-iam-auth/v0.1.2/kic-iam-auth_0.1.2_Linux_x86_64/kic-iam-auth && chmod +x kic-iam-auth && mv kic-iam-auth /usr/local/bin/

#-------------------------------------------
# STEP 3 : execute a binary
#-------------------------------------------
ENV KUBECONFIG  ""

ENTRYPOINT [ "/app/dashboard-metrics-scraper" ]

EXPOSE 8080

