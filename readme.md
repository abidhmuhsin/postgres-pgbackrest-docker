
- minio/ - local s3
- postgres/ - postgres+pgbackrest

## Test

clear existing build image i.e to cleanup all pgdata folder leftovers
> docker rm postgres-pgbackrest 

optionally cleanup s3 backup folder
> sudo rm -rf ./minio/_s3data/pgbackup/pgbackrest-test/

optionally cleanup pgdata folder
> sudo rm -rf ./postgres/_pgdata/*

rebuild container
> sudo docker-compose up --build

just restart container preserving all pgdata
> sudo docker-compose up

if no s3 bucket exists login to minio gui using default password and create bucket `pgbackup`

if stanza does not exist on backup folder, use create-stanza

create stanza example. [postgres-pgbackrest is the container name]
> docker exec -u postgres -it postgres-pgbackrest pgbackrest stanza-create --stanza="pgbackrest_test_stanza_name" --log-level-console=info

backup(default repo1)
> docker exec -u postgres -it postgres-pgbackrest pgbackrest backup --stanza="pgbackrest_test_stanza_name" --log-level-console=info

backup repo 2
> docker exec -u postgres -it postgres-pgbackrest pgbackrest backup --repo=2 --stanza="pgbackrest_test_stanza_name" --log-level-console=info

info
> docker exec -u postgres -it postgres-pgbackrest pgbackrest info --stanza="pgbackrest_test_stanza_name" --log-level-console=info

read postgres/readme.md for direct build and usage notes

view backup files

    repo1: check ip:9000 [https://172.77.0.45:9001] for minio gui or directly explore ./minio/s3data/pgbackup folder.
    repo2: check `docker exec -u postgres -it postgres-pgbackrest ls /pg_backup/data/pgbackrest/`
    logs: check `docker exec -u postgres -it postgres-pgbackrest ls /pg_backup/data/log/pgbackrest/`