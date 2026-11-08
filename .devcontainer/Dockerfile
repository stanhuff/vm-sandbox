FROM golang:1.25-bookworm
RUN useradd --create-home --uid 1000 --shell /bin/bash developer
ENV GOPATH=/home/developer/go
ENV PATH=/home/developer/go/bin:${PATH}
