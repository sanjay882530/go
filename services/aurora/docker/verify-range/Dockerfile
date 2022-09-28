FROM ubuntu:20.04

MAINTAINER Bartek Nowotarski <bartek@hcnet.org>

ARG HCNET_CORE_VERSION 
ENV HCNET_CORE_VERSION=${HCNET_CORE_VERSION:-*}
# to remove tzdata interactive flow
ENV DEBIAN_FRONTEND=noninteractive

ADD dependencies /
RUN ["chmod", "+x", "dependencies"]
RUN /dependencies

ADD hcnet-core.cfg /
ADD captive-core-pubnet.cfg /
RUN mkdir -p /cc
RUN mkdir -p /data

ADD start /
RUN ["chmod", "+x", "start"]

ENTRYPOINT ["/start"]
