FROM postgres:17

RUN apt-get update && apt-get install -y wget \
    && rm -rf /var/lib/apt/lists/*

# snag a binary of golang-migrate
RUN wget -nv https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \ 
    && tar -xzf migrate.linux-amd64.tar.gz \
    && rm migrate.linux-amd64.tar.gz \
    && cp migrate /bin/migrate

# copy migration script to be run at startup
COPY init.sh /docker-entrypoint-initdb.d/at_startup/init.sh
