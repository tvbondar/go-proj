FROM golang:1.25-bookworm
RUN apt-get update && apt-get install -y gcc
WORKDIR /app
COPY . .
ENV CGO_ENABLED=1
CMD ["go", "run", "-race", "l1-7.go"]