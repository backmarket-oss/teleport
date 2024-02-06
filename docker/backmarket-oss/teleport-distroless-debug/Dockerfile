FROM public.ecr.aws/gravitational/teleport-distroless-debug:14.3.6
LABEL org.opencontainers.image.source="https://github.com/backmarket-oss/teleport"
MAINTAINER Benoît Sauvère <benoit.sauvere@backmarket.com>

# Patch the teleport binary with the one from the fork
ADD .tmp/teleport/teleport /usr/local/bin/teleport