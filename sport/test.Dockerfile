FROM golang:1.16
ADD ./back /back
WORKDIR /back
ENTRYPOINT /bin/bash -c "make test"