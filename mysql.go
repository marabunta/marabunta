package marabunta

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var createTableStatements = []string{
	`CREATE TABLE IF NOT EXISTS states (
		id TINYINT(2) NOT NULL AUTO_INCREMENT,
		state VARCHAR(32) NOT NULL,
		PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin`,
	`CREATE TABLE IF NOT EXISTS payloads (
		id BINARY(16),
		name VARCHAR(255) NOT NULL,
		cdate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		payload JSON NOT NULL,
		PRIMARY KEY (id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin`,
	`CREATE TABLE IF NOT EXISTS tasks (
		id BINARY(16),
		cdate DATETIME NOT NULL,
		description VARCHAR(255),
		enabled TINYINT(1) unsigned NOT NULL DEFAULT 1,
		mdate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		name VARCHAR(64),
		payload_id binary(16),
		retried TINYINT unsigned,
		retries TINYINT unsigned DEFAULT 0,
		schedule VARCHAR(64) NOT NULL,
		sdate DATETIME,
		state_id TINYINT(2),
		target VARCHAR(255) NOT NULL,
		type TINYINT(2) unsigned NOT NULL,
		PRIMARY KEY (id),
		KEY (enabled),
		KEY (name),
		KEY (state_id),
		FOREIGN KEY (payload_id)
			REFERENCES payloads(id)
			ON DELETE CASCADE,
		FOREIGN KEY (state_id)
			REFERENCES states(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin`,
	`CREATE TABLE IF NOT EXISTS jobs (
		id BINARY(16) NOT NULL,
		task_id BINARY(16),
		state_id TINYINT,
		cdate DATETIME NOT NULL,
		mdate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (task_id)
			REFERENCES tasks(id)
			ON DELETE CASCADE,
		FOREIGN KEY (state_id)
			REFERENCES states(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin`,
	`CREATE TABLE IF NOT EXISTS messages (
		id BINARY(16) NOT NULL,
		job_id BINARY(16),
		cdate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		msg TEXT NOT NULL,
		PRIMARY KEY (id, job_id),
		FOREIGN KEY (job_id)
			REFERENCES jobs(id)
			ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin`,
}

func initMySQL(c *Config) (*sql.DB, error) {
	// test MySQL connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.MySQL.Username,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Database))
	if err != nil {
		return nil, fmt.Errorf("mysql: Error on initializing database connection: %s", err.Error())
	}

	// sql pool options
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(5)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("mysql: could not connect to the database: %s", err.Error())
	}

	// check if table exist if not create it
	if _, err := db.Exec("DESCRIBE tasks"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			if err := createTable(db); err != nil {
				return nil, err
			}
			return db, populateTable(db)
		}
		// Unknown error.
		return nil, fmt.Errorf("mysql: could not connect to the database: %v", err)
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func populateTable(db *sql.DB) error {
	// states
	var (
		sqlStr      = "INSERT INTO states(state) VALUES "
		sqlStrPstms []string
		states      = []string{"todo", "queued", "running", "done", "error"}
		vals        []interface{}
	)
	for _, row := range states {
		sqlStrPstms = append(sqlStrPstms, "(?)")
		vals = append(vals, row)
	}
	sqlStr = sqlStr + strings.Join(sqlStrPstms, ",")
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vals...)
	return err
}
