# Go + gRPC + Postgres

## Pre-requisites

### Important Postgresql commands

* Setup Postgres - `brew install postgresql`
* Start Postgres - `brew services start postgresql`
* Stop Postgres - `brew services stop postgresql`

### Configure Postgres
`psql postgres`

```
CREATE ROLE newUser WITH LOGIN PASSWORD ‘password’;
ALTER ROLE newUser CREATEDB;
```

Reconnect using the newUser
`psql postgres -U newuser`

Postgres default port is `:5432`

## Acknowledgments

* [Postgres Doc](https://www.sqlshack.com/setting-up-a-postgresql-database-on-mac/)
* [Tutorial series](https://www.youtube.com/watch?v=YudT0nHvkkE&list=PLrSqqHFS8XPYu-elDr1rjbfk0LMZkAA4X)

