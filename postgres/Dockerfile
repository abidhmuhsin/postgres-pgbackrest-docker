FROM postgres:14.2-alpine
# FROM postgres:13.6-alpine3.15
LABEL AUTHOR="Abidh Muhsin"

# Copy any necessary config files 
# ADD pg_hba.conf /var/lib/postgresql/data/
# ADD postgresql.conf /var/lib/postgresql/data/
# Add one or more *.sql, *.sql.gz, or *.sh scripts under /docker-entrypoint-initdb.d/ - ontime run for new db
# COPY init.sql /docker-entrypoint-initdb.d/

# install pgbackrest
RUN set -x -o pipefail \
    && apk update \
    && apk upgrade \
    && apk add \
    nano pgbackrest \
    && rm -rf /var/cache/apk/*

# if using pgbackrest to backup to local repo on a mountpoint, create necessary folders and give ownership to PGUSER
# using same postgres user for simplicity or setup fully as in https://blog.crunchydata.com/blog/secure-permissions-for-pgbackrest
RUN mkdir -p /pg_backup/data/pgbackrest
RUN mkdir -p /pg_backup/data/log/pgbackrest
RUN chown -R postgres:postgres /pg_backup/

# pgbackrest
# copy pgbackrestconfig
COPY conf/pgbackrest.conf /etc/pgbackrest/pgbackrest.conf


# use custom config file(archiving options or other setup), should be copied on init
# get the default config [docker run -i --rm postgres cat /usr/share/postgresql/postgresql.conf.sample > my-postgres.conf]
#x not work-- COPY conf/my-postgres.conf /var/lib/postgresql/data/postgresql.conf -- GETS OVERWRITTEN BY ENTRYPOINT CODE
COPY conf/my-postgres.conf /tmp/postgresql.conf
COPY init-scripts/0_update-conf.sh /docker-entrypoint-initdb.d/
# COPY init-scripts/init-alter-system-commands.sh /docker-entrypoint-initdb.d/


# ANY *.sql, *sh files under /docker-entrypoint-initdb.d/ will get executed during first init
# COPY init-scripts/init-user-db.sh /docker-entrypoint-initdb.d/
# COPY init-scripts/init.sql /docker-entrypoint-initdb.d/


RUN echo done