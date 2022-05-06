FROM ubuntu:focal

ARG REMOTE_USER
ARG REMOTE_PW

COPY "./init.sh" "/tmp"
COPY "./ssh-pub-key" "/tmp/ssh-pub-key"
RUN bash "/tmp/init.sh"

RUN mkdir "/docker"
COPY "./entrypoint.sh" "/docker"

EXPOSE 22

CMD ["/bin/bash", "/docker/entrypoint.sh"]
