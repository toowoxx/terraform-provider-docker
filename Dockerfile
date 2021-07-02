FROM ubuntu:21.04

RUN apt update -y && apt upgrade -y
RUN apt install -y golang bash make ca-certificates
