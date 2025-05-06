FROM mcr.microsoft.com/devcontainers/base:debian

# install nats cli

RUN curl -sf https://binaries.nats.dev/nats-io/natscli/nats@latest | sh \
    && mv nats /usr/local/bin/nats