FROM golang:1.18 as BUILDER

# caching stuff
COPY ./back/go.mod ./back/go.sum /back/
WORKDIR /back
RUN go mod download
COPY ./back /back
RUN mkdir -p /back/bin/static_files          \
    && make build                               \
    && find /back -type f -not -name "run_api"  \
        -and -not -name "api_conf_ver"          \
        -and -not -name "*.yml"                 \
        -and -not -name "*.sql"                 \
        -and -not -name "*.html"                \
        -delete



FROM centos:centos7
COPY --from=BUILDER /back /back
COPY ./lang /lang
RUN mkdir -p /back
EXPOSE 1580
WORKDIR /back/bin/                    
CMD ["./run_api"]
