#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    ALTER SYSTEM SET wal_level = replica;
    ALTER SYSTEM SET archive_mode = on;
    ALTER SYSTEM SET archive_command = 'pgbackrest --stanza=pgbackrest_test_stanza_name archive-push %p';
    ALTER SYSTEM SET max_wal_senders = 3;
    ALTER SYSTEM SET log_line_prefix = '';
    ALTER SYSTEM SET listen_addresses = '*'
EOSQL

#  ALTER SYSTEM SET archive_mode & archive_command was not reflecing in conf file but persistent after restarts and available when checking with 
# select name,setting from pg_settings where name like 'archive%' ;