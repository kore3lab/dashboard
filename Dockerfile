# #-------------------------------------------
# # STEP 1 : build executable binary
# #-------------------------------------------
FROM node:14.14.0-alpine3.12 as builder

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN npm install
RUN npm run build:frontend

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM node:14.14.0-alpine3.12

COPY --from=builder /usr/src/app /app/

ENV HOST 0.0.0.0
EXPOSE 3000

WORKDIR /app
ENTRYPOINT [ "npm", "run", "run" ]
