# ChirpStack PostgreSQL to SQLite script

This script connects to a ChirpStack PostgreSQL database, reads the data
and writes this into s SQLite database. It is intended to migrate the
ChirpStack Gateway OS (full) image from PostgreSQL to SQLite.

## Usage

### CLI

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

### ChirpStack Gateway OS

In the web-interface, under **System > Software**, click the **Upload package**
button and upload the `.ipk` package. This will automatically migrate the
PostgreSQL database to SQLite (leaving the PostgreSQL database as-is).

## Building from source

```text
# Build for current architecture
make build

# Build ChirpStack Gateway OS (ARMv7) package
make build-gateway-os
```

## Changelog

### v4.0.0

Initial release compatible with ChirpStack v4.9.0 database schema, ChirpStack
Gateway OS v4.5.x.

## License

This script is distributed under the MIT license. See also `LICENSE`.
