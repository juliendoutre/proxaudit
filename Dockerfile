FROM scratch
COPY proxaudit /
ENTRYPOINT ["/proxaudit"]
