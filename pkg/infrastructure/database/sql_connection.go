package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	dbConf "github.com/raymondsugiarto/reputation-be/config"

	"github.com/golang-migrate/migrate/v4/database"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type SQLConnection struct {
	config dbConf.Database
	dbConn *gorm.DB
	schema string
}

func NewSQLConnection(config dbConf.Database, schema string) (*SQLConnection, error) {
	db, err := connect(config, schema)
	if err != nil {
		return nil, err
	}

	return &SQLConnection{
		config: config,
		dbConn: db,
		schema: schema,
	}, nil
}

func connect(config dbConf.Database, schemaName string) (*gorm.DB, error) {
	dialect, err := getGormDialect(config, schemaName)
	if err != nil {
		return nil, err
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // Disable color
		},
	)

	dbConn, err := gorm.Open(dialect, &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	log.Printf("connected to db %+v", dbConn)

	//dbConn.SetLogger(&logger.GormLogger{})
	//dbConn.LogMode(true)

	return dbConn, err
}

func getGormDialect(config dbConf.Database, schema string) (gorm.Dialector, error) {
	sqlDB, err := getSqlDB(config, schema)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		return nil, err
	}
	var dialect gorm.Dialector
	if config.Adapter == "mysql" {
		dialect = mysql.New(mysql.Config{Conn: sqlDB})
	} else if config.Adapter == "postgres" {
		dialect = postgres.New(postgres.Config{Conn: sqlDB})
	}
	return dialect, nil
}

func getSqlDB(config dbConf.Database, schema string) (*sql.DB, error) {
	username := config.Username
	password := config.Password
	host := config.Host
	port := config.Port
	dbname := config.Dbname

	var dsn string
	if config.Adapter == "mysql" {
		dsn = username + ":" + password + "@(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
	} else if config.Adapter == "postgres" {
		dsn = "host=" + host + " user=" + username + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Asia/Jakarta"

		if schema != "" && schema != "public" {
			dsn += " search_path=" + schema
		}

	} else {
		// TODO: adapter not support
	}
	return sql.Open(config.Adapter, dsn)
}

func GetDatabaseDriverMigration(config dbConf.Database, schema string) (database.Driver, error) {
	sqlDB, err := getSqlDB(config, schema)
	if err != nil {
		return nil, err
	}
	var driver database.Driver
	if config.Adapter == "mysql" {
		driver, _ = mysqlMigrate.WithInstance(sqlDB, &mysqlMigrate.Config{})
	} else if config.Adapter == "postgres" {
		driver, _ = postgresMigrate.WithInstance(sqlDB, &postgresMigrate.Config{})
	}
	return driver, nil
}

func (d *SQLConnection) GetConn() *gorm.DB {
	return d.dbConn
}

var (
	DBConn *gorm.DB
)

// InitForCLI opens a GORM connection using the same config the
// production server reads, then exposes it via the package-level
// DBConn so CLI commands (like `app db seed`) can reuse it.
//
// Unlike the production server's initDatabase() helper — which
// log.Fatal()s on connection error — this returns the error so the
// caller can decide whether to abort or fall back. The CLI commands
// (cmd/db/db.go) choose to print the error and exit.
func InitForCLI() error {
	cfg := dbConf.GetConfig().Database.Main
	conn, err := NewSQLConnection(cfg, cfg.Schema)
	if err != nil {
		return err
	}
	DBConn = conn.GetConn()
	return nil
}
