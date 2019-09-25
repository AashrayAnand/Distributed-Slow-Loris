FROM golang:latest
LABEL maintainer="aashrayanand01@gmail.com"
RUN go get github.com/AashrayAnand/Distributed-Slow-Loris/worker
EXPOSE 3000
CMD ["worker"]
