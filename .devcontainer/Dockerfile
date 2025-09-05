FROM mcr.microsoft.com/vscode/devcontainers/go:1.24

RUN apt update
RUN apt install -y git cmake make g++ libssl-dev
RUN apt install -y postgresql
USER vscode

RUN go install github.com/vektra/mockery/v2@v2.50.0

CMD /bin/sh -c "while sleep 86000; do :; done"
