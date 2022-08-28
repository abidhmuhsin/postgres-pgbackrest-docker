
## To build and run the docker image 
> ./build.sh 

remove any prev builds and rebuild
>  docker rm abidh-postgres && ./build.sh

> docker build .-t abidh/postgres:latest


## Database Configuration
```sh
# get the default config from standard container
$ docker run -i --rm postgres cat /usr/share/postgresql/postgresql.conf.sample > my-postgres.conf
# customize the config
# run postgres with custom config or use Dockerfile for custom build
$ docker run -d --name some-postgres -v "$PWD/my-postgres.conf":/etc/postgresql/postgresql.conf -e POSTGRES_PASSWORD=mysecretpassword postgres -c 'config_file=/etc/postgresql/postgresql.conf'
```

## Setup pgbackrest

1. Define defaults for [global] and [stanza_name] in pgbackrest.conf
2. Copy pgbackrest.conf to /etc/pgbackrest/pgbackrest.conf
3. If using local mounts or directories for backup provide folder ownership/permissions for backup folder. If using s3, bucket should be precreated. Stanza folders will eb created automatically inside the bucket when u run create-stanza command.
4. Update archive_command/archive_mode and other archive config in postgres.conf. (See added lines towards last in sample config)
5. Use create-stanza stanza_name command for initializing backup folder first time. (While postgres is running )
6. To restore use restore command while postgres is stopped.


## Commands to run via pgbackrest docker exec
```sh
# drop into container shell
docker exec -it abidh-postgres bin/sh
# su postgres

# create stanza
docker exec -u postgres -it abidh-postgres pgbackrest stanza-create --stanza="pgbackrest_test_stanza_name" --log-level-console=info
# print backup info
docker exec -u postgres -it abidh-postgres pgbackrest info --stanza="pgbackrest_test_stanza_name"

# backup all databases, complete stanza, i.e all dbs ( first backup always type=full)
docker exec -u postgres -it abidh-postgres pgbackrest backup --stanza="pgbackrest_test_stanza_name" --log-level-console=info
# for redoing full-backup after making first backup
docker exec -u postgres -it abidh-postgres pgbackrest backup --stanza="pgbackrest_test_stanza_name" --type=full --log-level-console=detail

# restore full satnza(all dbs)
docker exec -u postgres -it abidh-postgres pgbackrest restore --stanza="pgbackrest_test_stanza_name" --delta --log-level-console=detail

# Restore only specified databases. (--db-include).NOTE: built-in databases (template0, template1, and postgres) are always restored unless specifically excluded.
# This feature allows only selected databases to be restored. Databases not specifically included will be restored as sparse, zeroed files to save space but still allow PostgreSQL to perform recovery. After recovery, the databases that were not included will not be accessible but can be removed with the drop database command.
--db-include=my_db

# Restore excluding the specified databases.
# Databases excluded will be restored as sparse, zeroed files to save space but still allow PostgreSQL to perform recovery.
# After recovery, those databases will not be accessible but can be removed with the drop database command.
# The --db-exclude option can be passed multiple times to specify more than one database to exclude. When used in combination with the --db-include option, --db-exclude will only apply to standard system databases (template0, template1, and postgres).
--db-exclude=postgres

# work with multiple backup repositories
--repo=1

# Type Option (--type)
# Recovery type. - The following recovery types are supported:
# default - recover to the end of the archive stream.
# immediate - recover only until the database becomes consistent. This option is only supported on PostgreSQL >= 9.4.
# lsn - recover to the LSN (Log Sequence Number) specified in --target. This option is only supported on PostgreSQL >= 10.
# name - recover the restore point specified in --target.
# xid - recover to the transaction id specified in --target.
# time - recover to the time specified in --target.
# preserve - preserve the existing recovery.conf file.
# standby - add standby_mode=on to recovery.conf file so cluster will start in standby mode.
# none - no recovery.conf file is written so PostgreSQL will attempt to achieve consistency using WAL segments present in pg_xlog/pg_wal. Provide the required WAL segments or use the archive-copy setting to include them with the backup.
--type=immediate
```

## Restore can also be done by editing commandproperty in docker-compose.yml

    # command: docker-entrypoint.sh postgres
    ## BACKUP LIVE :: works on running container. not necessary to run via command itself..can run outside container
    # command: "pgbackrest backup --stanza=pgbackrest_test_stanza_name"
    ## RESTORE ON LIVE CONTAINER:: postmaster.pid needs to be removed before running a restore on live/used data folder. Run required commands in sequence using `bash -c`
    # command: 
    #     - /bin/bash
    #     - -c
    #     - |
    #         rm -r /var/lib/postgresql/data/postmaster.pid
    #         pgbackrest restore --stanza=pgbackrest_test_stanza_name --delta --log-level-console=detail
    #         docker-entrypoint.sh postgres
    ## RESTORE FRESH:: only works on a clean pgdata folder,- uncomment below command with appropriate stanza name(--stanza=pgbackrest_test_stanza_name) to run restore before initializing the container.
    # command: pgbackrest restore --stanza=pgbackrest_BACKED_UP_stanza --repo1-path=/pgbackrest-test  --pg1-path=/var/lib/postgresql/data/ --delta --type=immediate --log-level-console=detail --recovery-option=recovery_target_action=promote
    # command: su - postgres -c 'postgres  -D /var/lib/postgresql/data --config-file=/var/lib/postgresql/data/postgresql.conf'
    

### View backup files

    repo1: check ip:9000 [https://172.77.0.45:9001] for minio gui or directly explore ./minio/s3data/pgbackup folder.
    repo2: execute `docker exec -u postgres -it postgres-pgbackrest ls /pg_backup/data/pgbackrest/`
    logs: execute `docker exec -u postgres -it postgres-pgbackrest ls /pg_backup/data/log/pgbackrest/`

## Links

>- https://severalnines.com/database-blog/validating-your-postgresql-backups-docker
>- https://www.enterprisedb.com/docs/supported-open-source/pgbackrest/08-multiple-repositories/
>- https://fatdba.com/2021/04/08/pgbackrest-a-reliable-backup-and-recovery-solution-to-your-postgresql-clusters/



 - https://vhutie.medium.com/custom-postgres-docker-image-with-predefined-database-and-tables-with-permanent-storage-2c3b44f92aad

- https://levelup.gitconnected.com/creating-and-filling-a-postgres-db-with-docker-compose-e1607f6f882f

- [If you are using a self-signed cert then you can use --repo1-s3-verify-ssl=n to disabled verification of the cert. There's no option to disable HTTPS for pgbackrest but minio can be configured to do HTTPS.]
https://blog.dbi-services.com/using-pgbackrest-to-backup-your-postgresql-instances-to-a-s3-compatible-storage/ [minio encrypted connection example]
https://docs.min.io/docs/how-to-secure-access-to-minio-server-with-tls.html

- https://hub.docker.com/_/postgres
- https://github.com/docker-library/postgres/blob/e8ebf74e50128123a8d0220b85e357ef2d73a7ec/14/alpine/Dockerfile


> ### Initialization scripts execution order in docker image.
> https://hub.docker.com/_/postgres/
>
> If you would like to do additional initialization in an image derived from this one, add one or more *.sql, *.sql.gz, or *.sh scripts under /docker-entrypoint-initdb.d (creating the directory if necessary). After the entrypoint calls initdb to create the default postgres user and database, it will run any *.sql files, run any executable *.sh scripts, and source any non-executable *.sh scripts found in that directory to do further initialization before starting the service.
Warning: scripts in /docker-entrypoint-initdb.d are only run if you start the container with a data directory that is empty; any pre-existing database will be left untouched on container startup. One common problem is that if one of your /docker-entrypoint-initdb.d scripts fails (which will cause the entrypoint script to exit) and your orchestrator restarts the container with the already initialized data directory, it will not continue on with your scripts.

### :: Notes (after restore)
-----------------
- ERROR: cannot execute CREATE TABLE in a read-only transaction.

    Since SELECT pg_is_in_recovery() is true you're connected to a read-only replica server in hot_standby mode. The replica configuration is in recovery.conf.
    You can't make it read/write except by promoting it to a master, at which point it will stop getting new changes from the old master server. See the PostgreSQL documentation on replication.

    First step is to check whether there is a 'recovery.conf' file in the data directory. If it exists, if you are sure you are on master (not slave) server, rename that file to 'recover.conf.backup'. Then, re-start postgresql server. It should allow you to write new records now.

    As you can see, recovery pauses after reaching the target.

    This is because you have hot_standby = on and you left recovery_target_action at its default value pause.

    You have to add the following in recovery.conf:

    recovery_target_action = promote (#--recovery-option=recovery_target_action=promote will work after next restart.)
    Alternatively, you can connect to the recovering server and complete recovery manually:
    SELECT pg_wal_replay_resume();
--------------------------------------


 - FATAL:  relation mapping file "base/13756/pg_filenode.map" contains invalid data

    Internal tables are not restored properly