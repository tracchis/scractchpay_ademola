FROM golang:1.15.3

WORKDIR /app

COPY . .

RUN apt-get update && \
    apt-get install -y \
    git \
    curl \
    jq


ENV PORT=8000
EXPOSE $PORT

CMD ["go", "build"]
CMD ["go", "run", "cmd/main.go", "cmd/config.go", "cmd/routes.go"]