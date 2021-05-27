#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM node:14.14.0-alpine3.12 as builder

ADD . /usr/src/app

WORKDIR /usr/src/app
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN npm install --no-optional
RUN npm run build:frontend

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM nginx:1.20.0-alpine

RUN rm -rf /usr/share/nginx/html/connect /usr/share/nginx/html/device /usr/share/nginx/html/event /usr/share/nginx/html/system /usr/share/nginx/html/tag /usr/share/nginx/html/test /usr/share/nginx/html/user /usr/share/nginx/html/auth /usr/share/nginx/html/emit-prop /usr/share/nginx/html/monitoring /usr/share/nginx/html/dashboard

COPY --from=builder /usr/src/app/dist /usr/share/nginx/html
COPY ./scripts/docker/includes/default.conf /etc/nginx/conf.d/default.conf
COPY ./scripts/docker/includes/nginx.conf /etc/nginx/nginx.conf

ENV HOST 0.0.0.0
ENV BACKEND_PORT 3001

ENTRYPOINT [ "nginx", "-g", "daemon off;" ]
