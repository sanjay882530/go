FROM golang:1.19.1 as build

ADD . /src/soroban-rpc
WORKDIR /src/soroban-rpc
RUN go build -o /bin/soroban-rpc ./exp/services/soroban-rpc


FROM ubuntu:20.04
ARG HCNET_CORE_VERSION
ENV HCNET_CORE_VERSION=${HCNET_CORE_VERSION:-*}
ENV HCNET_CORE_BINARY_PATH /usr/bin/hcnet-core
ENV DEBIAN_FRONTEND=noninteractive

# ca-certificates are required to make tls connections
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl wget gnupg apt-utils
RUN wget -qO - https://apt.hcnet.org/SDF.asc | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=true apt-key add -
RUN echo "deb https://apt.hcnet.org focal stable" >/etc/apt/sources.list.d/SDF.list
RUN echo "deb https://apt.hcnet.org focal unstable" >/etc/apt/sources.list.d/SDF-unstable.list
RUN apt-get update && apt-get install -y hcnet-core=${HCNET_CORE_VERSION}
RUN apt-get clean

COPY --from=build /bin/soroban-rpc /app/
ENTRYPOINT ["/app/soroban-rpc"]
