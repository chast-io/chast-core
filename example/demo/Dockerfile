FROM ubuntu:22.10
MAINTAINER chast.io <Raphael Jenni>

# ==================== IMPORTANT =================================
# Linux Kernel > v4.18.0 is required
# ================================================================

RUN apt-get update && apt-get install -y unionfs-fuse mergerfs libfuse3-dev fuse-overlayfs tree

RUN apt-get install -y openjdk-17-jdk openjdk-17-jre

# Run go build main.go in the cli folder first
COPY main /usr/local/bin/chast

COPY recipes /recipes

RUN echo 'echo "\e[1;34mImportant Note: CHAST requires Linux Kernel > v4.18.0 to run properly\e[0m"' > /etc/profile.d/welcome.sh

ENTRYPOINT ["/bin/bash", "-l"]
