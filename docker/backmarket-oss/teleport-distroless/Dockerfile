FROM public.ecr.aws/gravitational/teleport-distroless:16.1.0
ABEL org.opencontainers.image.source="https://github.com/backmarket-oss/teleport"
MAINTAINER Benoît Sauvère <benoit.sauvere@backmarket.com>

# Patch the teleport binary with the one from the fork
ADD .tmp/teleport/teleport /usr/local/bin/teleport