FROM mcr.microsoft.com/vscode/devcontainers/go:1.25

RUN apt update
RUN apt install -y git cmake make g++ libssl-dev
USER vscode

CMD /bin/sh -c "while sleep 86000; do :; done"
