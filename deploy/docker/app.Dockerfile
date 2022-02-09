FROM go-base:latest

WORKDIR /app

RUN mkdir -p /priv
COPY ./tmp/docker-db-within-containers.yml /priv/docker-postgres.yml
COPY ./tmp/priv.yml /priv
ENV APP_ENV_DB='/priv/docker-postgres.yml'
ENV SECRETS_FILE_PATH='/priv/priv.yml'

COPY . ./

RUN go build -o ./bin/exec ./main.go

EXPOSE 8000

CMD ["/app/bin/serv/http"]
