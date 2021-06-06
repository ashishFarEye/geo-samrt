package postgres

import (
	"crypto/tls"
	"fmt"

	"github.com/go-pg/migrations/v7"
	"go.uber.org/zap"

	"github.com/go-pg/pg/v9"
)

func New(url string, user string, pwd string, name string) *pg.DB {
	tls := tls.Config{
		InsecureSkipVerify: true,
	}

	db := pg.Connect(&pg.Options{
		Addr:            url,
		User:            user,
		Password:        pwd,
		Database:        name,
		ApplicationName: "maps",
		TLSConfig:       &tls,
	})
	return db
}

func Migrate(log *zap.Logger, db *pg.DB) {
	c := migrations.DefaultCollection
	if _, errr := migrations.Version(db); errr != nil {
		c.Run(db, "init")
	}
	if err := c.DiscoverSQLMigrations("db.migrations"); err != nil {
		log.Error("---Error Discovering Migrations---", zap.Error(err))
	}
	oldVersion, newVersion, err := c.Run(db)
	if err != nil {
		//log.Fatal("----Error While Running Migrations----", zap.Error(err))
	}
	if newVersion != oldVersion {
		log.Info(fmt.Sprintf("migrated from version %d to  %d", oldVersion, newVersion))
	} else {
		log.Info(fmt.Sprintf("version is  %d", oldVersion))
	}
}
