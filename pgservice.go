// Support PGSERVICE and PGSERVICEFILE environment variables for PostgreSQL.
package pgservice

import (
  "fmt"
  "os"

  "github.com/go-ini/ini"
)

// Valid parameters supported in reading `PGSERVICEFILE`.
var validParams = map[string]string{
  "host":        "PGHOST",
  "port":        "PGPORT",
  "dbname":      "PGDATABASE",
  "user":        "PGUSER",
  "password":    "PGPASSWORD",
  "sslmode":     "PGSSLMODE",
  "sslcert":     "PGSSLCERT",
  "sslkey":      "PGSSLKEY",
  "sslrootcert": "PGSSLROOTCERT",
}

// Apply config options specified in `PGSERVICE` section of `PGSERVICEFILE`
// (defaults to `~/.pg_service.conf`) to environment variables. This make them
// available to [`pq`](https://github.com/lib/pq). It also removes `PGSERVICE`
// and `PGSERVICEFILE` from environment variables to prevent `pq` panic.
// Read more about [The Connection Service File](https://www.postgresql.org/docs/current/static/libpq-pgservice.html).
func Apply() error {
  params, err := ReadFile()
  if err != nil {
    return err
  }

  for k, v := range params {
    if _, set := os.LookupEnv(validParams[k]); !set {
      os.Setenv(validParams[k], v)
    }
  }

  unsetenv()
  return nil
}

// ReadFile reads options specified in PGSERVICE section of PGSERVICEFILE
// (defaults to `~/.pg_service.conf`) and returns them as a map.
func ReadFile() (map[string]string, error) {
  result := make(map[string]string)

  pgservice, ok := os.LookupEnv("PGSERVICE")
  if !ok {
    return result, fmt.Errorf("PGSERVICE environment variable is not set")
  }

  pgservicefile, ok := os.LookupEnv("PGSERVICEFILE")
  if !ok {
    pgservicefile = os.ExpandEnv("${HOME}/.pg_service.conf")
  }

  cfg, err := ini.Load(pgservicefile)
  if err != nil {
    return result, fmt.Errorf("error loading PGSERVICEFILE at '%s'", pgservicefile)
  }

  cfg.BlockMode = false

  section, err := cfg.GetSection(pgservice)
  if err != nil {
    return result, err
  }

  for key := range validParams {
    if value, err := section.GetKey(key); err == nil {
      result[key] = value.String()
    }
  }

  return result, nil
}

// Prevents pq panic.
func unsetenv() {
  os.Unsetenv("PGSERVICE")
  os.Unsetenv("PGSERVICEFILE")
}
