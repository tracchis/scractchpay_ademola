FROM golang:1.15.3

WORKDIR /app

RUN apt-get update && \
    apt-get install -y \
    git \
    curl \
    jq

CMD ["go", "build"]
CMD ["go", "run", "cmd/main.go", "cmd/config.go", "cmd/routes.go"]
