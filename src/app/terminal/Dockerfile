#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.15.3-alpine3.12 as builder

# gcc
RUN apk add --no-cache build-base

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags kore-board.terminal -o kore-board.terminal -v main.go

#-------------------------------------------
# STEP 2 : terminal
#-------------------------------------------
FROM ubuntu:18.04

USER root

ARG ARCH=amd64

RUN rm -f /bin/sh && ln -s /bin/bash /bin/sh
ENV KUBECTL_VERSION v1.20.2

COPY --from=builder /usr/src/app/kore-board.terminal /app/

RUN apt-get update && \
  apt-get install -y --no-install-recommends curl ca-certificates jq iproute2 less bash-completion unzip sysstat acl net-tools iputils-ping telnet dnsutils wget vim git && \
  curl -sLf https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/${ARCH}/kubectl > /usr/bin/kubectl && chmod +x /usr/bin/kubectl && \
  git clone --branch master --depth 1 https://github.com/ahmetb/kubectl-aliases /opt/kubectl-aliases && chmod -R 755 /opt/kubectl-aliases && \
  curl -L https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash && \
  chmod 500 /usr/bin/nohup && \
  DEBIAN_FRONTEND=noninteractive apt-get --purge remove -y git && \
  DEBIAN_FRONTEND=noninteractive apt-get autoremove -y && \
  DEBIAN_FRONTEND=noninteractive apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
  chmod -R 755 /tmp && mkdir -p /opt/k3webterminal

COPY /scripts/vimrc.local /etc/vim
COPY /scripts/k3webterminal.sh /opt/k3webterminal/
COPY /scripts/config-namespace.sh /opt/k3webterminal/
COPY /scripts/config-filesystem.sh /opt/k3webterminal/
COPY /scripts/entryshell.sh /opt/k3webterminal/
RUN chmod -R 700 /opt/k3webterminal

ENV SESSION_STORAGE_SIZE=10M
#ENV WELCOME_BANNER="Welcome to K3WebTerminal for Kubernetes, try kubectl --help."
ENV KUBECTL_INSECURE_SKIP_TLS_VERIFY=true
## end install WebTerminal addon

# Declare the port on which the webserver will be exposed.
EXPOSE 3003

ENV KUBECONFIG ""

ENTRYPOINT [ "/app/kore-board.terminal" ]