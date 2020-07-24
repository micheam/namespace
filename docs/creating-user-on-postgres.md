# Creating User on postgres.

    export PSQL_PASSWD="<passwd>"
    psql -h localhost -U postgres -c "CREATE USER ${USER} WITH ENCRYPTED PASSWORD '${PSQL_PASSWD}';"
    psql -h localhost -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE ns TO ${USER};"

