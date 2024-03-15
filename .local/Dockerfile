FROM ubuntu:latest

# Set environment variables
ARG TARGETARCH
ENV CLIPSE_VERSION=0.0.6

RUN apt-get update && \
    apt-get install -y wget && \
    rm -rf /var/lib/apt/lists/*

RUN if [ "$TARGETARCH" = "arm64" ]; then \
        wget -c https://github.com/savedra1/clipse/releases/download/v${CLIPSE_VERSION}/clipse_${CLIPSE_VERSION}_linux_arm64.tar.gz -O - | tar -xz ; \
    elif [ "$TARGETARCH" = "amd64" ]; then \
        wget -c https://github.com/savedra1/clipse/releases/download/v${CLIPSE_VERSION}/clipse_${CLIPSE_VERSION}_linux_amd64.tar.gz -O - | tar -xz ; \
    elif [ "$TARGETARCH" = "836" ]; then \
        wget -c https://github.com/savedra1/clipse/releases/download/v${CLIPSE_VERSION}/clipse_${CLIPSE_VERSION}_linux_836.tar.gz -O - | tar -xz ; \
    elif [ "$TARGETARCH" = "darwin_arm64" ]; then \
        wget -c https://github.com/savedra1/clipse/releases/download/v${CLIPSE_VERSION}/clipse_${CLIPSE_VERSION}_darwin_arm64.tar.gz -O - | tar -xz ; \
    elif [ "$TARGETARCH" = "darwin_amd64" ]; then \
        wget -c https://github.com/savedra1/clipse/releases/download/v${CLIPSE_VERSION}/clipse_${CLIPSE_VERSION}_darwin_amd64.tar.gz -O - | tar -xz ; \
    else \
        echo "Unsupported architecture" ; \
    fi \
    mkdir .config

# Set the entrypoint
CMD ["/bin/bash"]