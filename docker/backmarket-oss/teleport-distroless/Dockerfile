FROM public.ecr.aws/gravitational/teleport-distroless:17.2.1
LABEL org.opencontainers.image.source="https://github.com/backmarket-oss/teleport"
LABEL org.opencontainers.image.authors="ops@backmarket.com"

# Patch the teleport binary with the one from the fork
ADD .tmp/teleport/teleport /usr/local/bin/teleport
