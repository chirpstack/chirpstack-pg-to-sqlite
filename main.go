package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	pgDB       *sqlx.DB
	slDB       *sqlx.DB
	sqlitePath string
	pgDSN      string
)

func init() {
	flag.StringVar(&sqlitePath, "sqlite-path", "chirpstack.sqlite", "Path to SQLite directory")
	flag.StringVar(&pgDSN, "postgres-dsn", "postgres://chirpstack:chirpstack@localhost/chirpstack?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()
}

func main() {
	log.Println("Opening databases")
	pgDB = getPostgresClient(pgDSN)
	slDB = getSqliteClient(sqlitePath)

	log.Println("Start migration")
	migrateTableFn("user", map[int]func(interface{}) interface{}{
		0: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
	})
	migrateTableFn("tenant", map[int]func(interface{}) interface{}{
		0:  fixUuid,
		1:  fixDateTime,
		2:  fixDateTime,
		10: fixKeyValue,
	})
	migrateTableFn("api_key", map[int]func(interface{}) interface{}{
		0: fixUuid,
		1: fixDateTime,
		4: fixUuid,
	})
	migrateTableFn("tenant_user", map[int]func(interface{}) interface{}{
		0: fixUuid,
		1: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
	})
	migrateTableFn("device_profile_template", map[int]func(interface{}) interface{}{
		0: fixUuid,
		1: fixDateTime,
		2: fixDateTime,
	})
	migrateTableFn("application", map[int]func(interface{}) interface{}{
		0: fixUuid,
		1: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
	})
	migrateTableFn("application_integration", map[int]func(interface{}) interface{}{
		0: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
	})
	migrateTableFn("gateway", map[int]func(interface{}) interface{}{
		1: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
		4: fixDateTime,
	})
	migrateTableFn("device_profile", map[int]func(interface{}) interface{}{
		0: fixUuid,
		1: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
	})
	migrateTableFn("device", map[int]func(interface{}) interface{}{
		1: fixUuid,
		2: fixUuid,
		3: fixDateTime,
		4: fixDateTime,
		5: fixDateTime,
	})
	migrateTableFn("device_keys", map[int]func(interface{}) interface{}{
		1: fixDateTime,
		2: fixDateTime,
		5: fixDevNonces,
	})
	migrateTableFn("device_queue_item", map[int]func(interface{}) interface{}{
		0:  fixUuid,
		2:  fixDateTime,
		8:  fixDateTime,
		10: fixDateTime,
	})
	migrateTableFn("multicast_group", map[int]func(interface{}) interface{}{
		0: fixUuid,
		1: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
	})
	migrateTableFn("multicast_group_device", map[int]func(interface{}) interface{}{
		0: fixUuid,
		2: fixDateTime,
	})
	migrateTableFn("multicast_group_gateway", map[int]func(interface{}) interface{}{
		0: fixUuid,
		2: fixDateTime,
	})
	migrateTableFn("multicast_group_queue_item", map[int]func(interface{}) interface{}{
		0: fixUuid,
		1: fixDateTime,
		2: fixDateTime,
		3: fixUuid,
	})
	migrateTableFn("relay_device", map[int]func(interface{}) interface{}{
		2: fixDateTime,
	})
	migrateTableFn("relay_gateway", map[int]func(interface{}) interface{}{
		0: fixUuid,
		2: fixDateTime,
		3: fixDateTime,
		4: fixDateTime,
	})
}

func fixDateTime(dt interface{}) interface{} {
	if dt == nil {
		return dt
	}

	dtCasted := dt.(time.Time)
	return dtCasted.UTC().Format("2006-01-02 15:04:05")
}

func fixUuid(id interface{}) interface{} {
	if id == nil {
		return id
	}

	var idUuid uuid.UUID
	err := idUuid.UnmarshalText(id.([]byte))
	if err != nil {
		log.Fatal("Uuid error", err)
	}

	return idUuid.String()
}

func fixKeyValue(kv interface{}) interface{} {
	if kv == nil {
		return kv
	}

	return string(kv.([]byte))
}

func fixDevNonces(v interface{}) interface{} {
	if v == nil {
		return v
	}

	str := string(v.([]byte))
	str = strings.Replace(str, "{", "[", 1)
	str = strings.Replace(str, "}", "]", 1)

	return str
}

func getPostgresClient(dsn string) *sqlx.DB {
	d, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Open PostgreSQL connection error", err)
	}

	return d
}

func getSqliteClient(path string) *sqlx.DB {
	d, err := sqlx.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal("Open SQLite error", err)
	}

	_ = d.MustExec("delete from tenant")
	_ = d.MustExec("delete from user")

	return d
}

func migrateTableFn(tableName string, modifiers map[int]func(interface{}) interface{}) {
	log.Println("Migrating", tableName, "table")

	rows, err := pgDB.Queryx(fmt.Sprintf("SELECT * from \"%s\"", tableName))
	if err != nil {
		log.Fatal("Query table error", err)
	}

	for rows.Next() {
		cols, err := rows.SliceScan()
		if err != nil {
			log.Fatal("Slice scan error", err)
		}

		if modifiers != nil {
			for i, f := range modifiers {
				cols[i] = f(cols[i])
			}
		}

		fields := []string{}

		for i := 1; i <= len(cols); i++ {
			fields = append(fields, fmt.Sprintf("$%d", i))
		}

		_, err = slDB.Exec(fmt.Sprintf("INSERT INTO \"%s\" values (%s)", tableName, strings.Join(fields, ", ")), cols...)
		if err != nil {
			log.Fatal("Insert error", err)
		}
	}

}
