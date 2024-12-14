FROM scratch

LABEL org.opencontainers.image.authors Julien Doutre <jul.doutre@gmail.com>
LABEL org.opencontainers.image.title proxaudit
LABEL org.opencontainers.image.url https://github.com/juliendoutre/proxaudit
LABEL org.opencontainers.image.documentation https://github.com/juliendoutre/proxaudit
LABEL org.opencontainers.image.source https://github.com/juliendoutre/proxaudit
LABEL org.opencontainers.image.licenses MIT

COPY proxaudit /
ENTRYPOINT ["/proxaudit"]
