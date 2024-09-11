package dal

import (
	"fmt"

	"github.com/asciifaceman/emri/pkg/dal/models"
	"github.com/asciifaceman/emri/pkg/global"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() (*PG, error) {
	p := &PG{
		l: zap.S().Named("postgres"),
	}
	return p, p.Connect()
}

type PG struct {
	l  *zap.SugaredLogger
	db *gorm.DB
}

func (p *PG) Connect() error {
	// TODO: connection pooling?
	dsn := global.C().PostgresConfig.DSN()
	dbn := global.C().PostgresConfig.Database

	p.l.Debugw("opening connection to db", "server", global.C().PostgresConfig.Hostname, "database", dbn, "username", global.C().PostgresConfig.Runtime.Username)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewLogger("gorm", dbn),
	})
	if err != nil {
		return err
	}
	p.db = db

	p.l = p.l.With("database", p.db.Name(), "sever", global.C().PostgresConfig.Hostname)
	p.l.Debugw("connection established with database")

	return nil
}

func (p *PG) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (p *PG) Migrate(obj ...interface{}) error {
	return p.db.AutoMigrate(obj...)
}

func (p *PG) Raw(format string, args ...interface{}) *gorm.DB {
	return p.db.Raw(fmt.Sprintf(format, args...))
}

func (p *PG) Exec(format string, args ...interface{}) *gorm.DB {
	return p.db.Exec(fmt.Sprintf(format, args...))
}

func (p *PG) CreateInstanceObject(domain *models.InstanceObject) error {
	p.l.Debugw("inserting domain", "domain", domain.Domain)

	if err := p.db.Create(domain); err.Error != nil {
		return err.Error
	}

	return nil
}
