FROM alpine:latest
WORKDIR /
ADD ./out/gorelay /
ADD ./.config.toml /
CMD ["./gorelay"]