# go-pgservice

Read [Connection Service File](https://www.postgresql.org/docs/current/static/libpq-pgservice.html) of PostgreSQL.

This library reads and applies options specified in a section of service file optionally overriden by `PGSERVICE` and `PGSERVICEFILE` environment variables.

### Install

```
go get github.com/behrang/go-pgservice
```

### Usage

Call `Apply` function to apply options in `service` section of `file` which can be overriden by `PGSERVICE` and `PGSERVICEFILE` environment variables. This make them available to [`pq`](https://github.com/lib/pq) as environment variables.

```go
package main

import (
  "database/sql"
  "log"

  "github.com/behrang/go-pgservice"
  _ "github.com/lib/pq"
)

func main() {
  err := pgservice.Apply("service", ".pg_service.conf")
  if err != nil {
    log.Fatal(err)
  }

  db, err := sql.Open("postgres", "")
  if err != nil {
    log.Fatal(err)
  }
  if err := db.Ping(); err != nil {
    log.Fatal(err)
  }
}
```

Now running your app with something like `PGSERVICE=mydb app` will run your app and connect to Postgre database specified in `[mydb]` section of `~/.pg_service.conf`.

### Supported Options In Connection Service File

```
host
port
dbname
user
password
sslmode
sslcert
sslkey
sslrootcert
```

### License

[MIT](LICENSE)
