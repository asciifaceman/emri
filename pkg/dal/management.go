package dal

/*
func Initialize() error {
	schema := global.C().PostgresConfig.Schema
	rtUser := global.C().PostgresConfig.Runtime.Username

	il := zap.S().Named("initializer")
	il.Info("Beginning database initialization")
	p, err := New(true)
	if err != nil {
		return err
	}

	if derr := p.DropDB(schema); derr != nil {
		return derr
	}

	il.Info("creating database")
	if cerr := p.CreateDB(schema, rtUser); cerr != nil {
		return cerr
	}

	il.Info("reconnecting under database")
	if err := p.Close(); err != nil {
		return err
	}

	p, err = New(false)
	if err != nil {
		return err
	}

	il.Info("setting up database")
	if err := p.Migrate(models.Migrate...); err != nil {
		return err
	}

	return nil

}

func (p *PG) DropDB(name string) error {
	if !p.manage {
		return errors.New("not connected as a management interface")
	}

	p.l.Warn("dropping database")

	stmt := "SELECT * FROM pg_database WHERE datname = '%s';"
	rs := p.Raw(stmt, name)
	if rs.Error != nil {
		return rs.Error
	}

	stmt = "DROP DATABASE %s;"
	var recv = make(map[string]interface{})
	if rs.Find(recv); len(recv) > 0 {
		if cs := p.Exec(stmt, name); cs.Error != nil {
			return cs.Error
		}
	} else {
		p.l.Errorw("database doesn't exist", "database", name)
	}

	return nil
}

func (p *PG) CreateDB(name string, owner string) error {
	if !p.manage {
		return errors.New("not connected as a management interface")
	}

	p.l.Info("creating database")

	stmt := "CREATE DATABASE %s WITH OWNER = %s ENCODING = 'UTF8' IS_TEMPLATE = False;"

	if cs := p.Exec(stmt, name, owner); cs.Error != nil {
		return cs.Error
	}

	stmt = "GRANT ALL PRIVILEGES ON DATABASE %s TO %s"

	if gs := p.Exec(stmt, name, owner); gs.Error != nil {
		return gs.Error
	}

	tx := p.db.Begin()
	tx.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON SCHEMA 'emri_dev.public' TO %s;", owner))
	if gs := tx.Commit(); gs.Error != nil {
		return gs.Error
	}

	return p.Close()
}

/*
Initialize the database from scratch, including dropping the db if it exists.

func ManagedInitialize() error {
	il := zap.S().Named("initializer")

	admin, err := New(true)
	if err != nil {
		return err
	}

	il.Info("Beginning initialization of database.")

	il.Info("Dropping db if exists...")
	name := global.C().PostgresConfig.Schema
	rs := admin.db.Raw(fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", name))
	if rs.Error != nil {
		return rs.Error
	}

	il.Info("Creating db...")
	stmt := "CREATE DATABASE %s WITH OWNER = postgres ENCODING = 'UTF8' IS_TEMPLATE = False;"

	var recv = make(map[string]interface{})
	if rs.Find(recv); len(recv) == 0 {
		if cs := admin.db.Exec(fmt.Sprintf(stmt, name)); cs.Error != nil {
			return cs.Error
		}
	} else {
		il.Errorw("database already exists during initialization after dropping", "database", name)
	}

	err = admin.Close()
	if err != nil {
		il.Error("failed to close management connection to database")
	}

	p, err := New(false)
	if err != nil {
		return err
	}

	err = p.Migrate(models.Migrate...)
	if err != nil {
		return err
	}

	return p.Close()
}
*/
