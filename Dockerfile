#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM node:14.14.0-alpine3.12 as builder

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN npm install --no-optional
RUN npm run build:frontend

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM node:14.14.0-alpine3.12

COPY --from=builder /usr/src/app/.nuxt /app/.nuxt
COPY --from=builder /usr/src/app/nuxt.config.js /app/nuxt.config.js
COPY --from=builder /usr/src/app/package.json /app/package.json

WORKDIR /app
RUN npm install --only=production  --no-optional

ENV HOST 0.0.0.0
ENV BACKEND_PORT 3001
ENV KIALI_PORT 20001

ENTRYPOINT [ "npm", "run", "run" ]
EXPOSE 3000