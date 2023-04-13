# stage 0 : npm stuff

FROM node:14 as npm_build

# caching stuff

ADD ui2/package*.json ui2/Makefile /ui2/
WORKDIR /ui2

RUN make install-prod

# building stuff

ADD ui2 /ui2

# will output stuff into /ui2/build
RUN make build

## stage 1 : nginx stuff
#
#FROM nginx:latest
#
#ARG nginx
#
#COPY --from=npm_build /ui2/build /usr/share/nginx/html
#
## copy config
#
#COPY ${nginx} /etc/nginx/conf.d/default.conf