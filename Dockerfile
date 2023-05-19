FROM ubuntu:latest
LABEL authors="eep"

ENTRYPOINT ["top", "-b"]