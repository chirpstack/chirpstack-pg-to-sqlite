# ChirpStack PostgreSQL to SQLite script

This script connects to a ChirpStack PostgreSQL database, reads the data
and writes this into s SQLite database. It is intended to migrate the
ChirpStack Gateway OS (full) image from PostgreSQL to SQLite.

## Usage

```text
Usage of ./chirpstack-pg-to-sqlite:
  -postgres-dsn string
        PostgreSQL DSN (default "postgres://chirpstack:chirpstack@localhost/chirpstack?sslmode=disable")
  -sqlite-path string
        Path to SQLite directory (default "chirpstack.sqlite")
```

**Important:**

* The target SQLite database needs to be already initialized with the proper
  schema. You can use the `chirpstack.empty.sqlite` as a skeleton.
* This script will remove all data from the target SQLite database.

## License

This script is distributed under the MIT license. See also `LICENSE`.
