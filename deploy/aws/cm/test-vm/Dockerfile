FROM ubuntu:latest

ARG PRIV_REMOTE_USER
ARG PRIV_REMOTE_PW

COPY "./init.sh" "/tmp"
COPY "./ssh-pub-key" "/tmp/ssh-pub-key"
RUN bash "/tmp/init.sh"

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D", "-o", "ListenAddress=0.0.0.0"]
