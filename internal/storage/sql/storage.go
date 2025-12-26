package sqlstorage

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"                                              //nolint:depguard
	_ "github.com/jackc/pgx/stdlib"                                       //nolint:depguard
	"github.com/jmoiron/sqlx"                                             //nolint:depguard
	"github.com/pressly/goose/v3"                                         //nolint:depguard
	"github.com/pressly/goose/v3/database"                                //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/storage/sql/migrations" //nolint:depguard
)

type Storage struct {
	db  *sqlx.DB
	c   config.StorageConfig
	err error
	ctx *context.Context
}

func New(ctx *context.Context, c config.StorageConfig) common.StorageDriverInterface {
	s := &Storage{c: c, ctx: ctx}
	s.err = s.Connect(*ctx)
	return s
}

func (s *Storage) Connect(ctx context.Context) error {
	s.db, s.err = sqlx.ConnectContext(ctx, "pgx", s.c.Dsn)
	if s.err != nil {
		return s.err
	}
	s.err = s.db.PingContext(ctx)
	return s.err
}

func (s *Storage) Close() error {
	s.err = s.db.Close()
	return s.err
}

func (s *Storage) Add(ipSubnet common.IPSubnet) (common.IPSubnet, error) {
	if ipSubnet.ID.(string) == "" {
		ipSubnet.ID = uuid.New().String()
	}
	sql := `INSERT INTO IPSubnets("id","title","date_time","duration","description","user","notify_time") 
VALUES(:id, :title, :date_time, :duration, :description, :user, :notify_time)`
	_, err := s.db.NamedExecContext(*s.ctx, sql, ipSubnet)
	if err != nil {
		return ipSubnet, err
	}
	return ipSubnet, err
}

func (s *Storage) Update(ipSubnet common.IPSubnet) error {
	sql := `UPDATE IPSubnets SET "title" = :title,"date_time" = :date_time,"duration" = :duration,
                  "description" = :description,"user" = :user,
                  "notify_time" = :notify_time WHERE id = :id`
	_, err := s.db.NamedExecContext(*s.ctx, sql, ipSubnet)
	if err != nil {
		return err
	}
	return err
}

func (s *Storage) Delete(id interface{}) error {
	sql := `DELETE FROM IPSubnets WHERE id = $1`
	_, err := s.db.ExecContext(*s.ctx, sql, id)
	return err
}

func (s *Storage) GetByID(id interface{}) (common.IPSubnet, error) {
	IPSubnet := common.IPSubnet{}
	sql := `SELECT "id","title","date_time","duration","description","user","notify_time" FROM IPSubnets WHERE id = $1`
	err := s.db.GetContext(*s.ctx, &IPSubnet, sql, id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return IPSubnet, common.ErrIPSubnetNotFound
	}
	return IPSubnet, err
}

func (s *Storage) List() ([]*common.IPSubnet, error) {
	IPSubnet := make([]*common.IPSubnet, 0)
	sql := `SELECT "id","title","date_time","duration","description","user","notify_time" FROM IPSubnets`
	err := s.db.SelectContext(*s.ctx, &IPSubnet, sql)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return IPSubnet, common.ErrIPSubnetNotFound
	}
	return IPSubnet, err
}

func (s *Storage) PrepareStorage(log common.LoggerInterface) error {
	provider, err := goose.NewProvider(database.DialectPostgres, s.db.DB, migrations.Embed)
	if err != nil {
		log.Error("init goose", "error", err)
		return err
	}
	sources := provider.ListSources()
	for _, s := range sources {
		log.Info("Migration item", "type", s.Type, "version", s.Version, "path", filepath.Base(s.Path))
	}

	stats, err := provider.Status(*s.ctx)
	if err != nil {
		log.Error("status", "error", err)
		return err
	}
	for _, s := range stats {
		log.Info("Migrate status", "type", s.Source.Type, "version", s.Source.Version, "duration", s.State)
	}
	results, err := provider.Up(*s.ctx)
	if err != nil {
		log.Error("up", "error", err)
		return err
	}
	for _, r := range results {
		log.Info("Migrate done", "type", r.Source.Type, "version", r.Source.Version, "duration", (r.Duration).String())
	}

	return nil
}

func (s *Storage) IsOverlapping(_ *common.IPSubnet) (bool, error) {
	count := 0
	var err error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
