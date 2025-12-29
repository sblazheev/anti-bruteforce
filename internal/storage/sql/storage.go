package sqlstorage

import (
	"context"
	"fmt"
	"path/filepath"

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

func (s *Storage) Add(jar string, ipSubnet common.IPSubnet) (*common.IPSubnet, error) {
	sql := `INSERT INTO %s("sub_net") VALUES(:sub_net) ON CONFLICT DO NOTHING`
	sql = fmt.Sprintf(sql, jar)
	_, err := s.db.NamedExecContext(*s.ctx, sql, ipSubnet)
	if err != nil {
		return &ipSubnet, err
	}
	return &ipSubnet, err
}

func (s *Storage) Update(jar string, ipSubnet common.IPSubnet) error {
	sql := `UPDATE %s SET "sub_net" = :sub_net,"date_create" = :date_create WHERE sub_net = :sub_net`
	sql = fmt.Sprintf(sql, jar)
	_, err := s.db.NamedExecContext(*s.ctx, sql, ipSubnet)
	if err != nil {
		return err
	}
	return err
}

func (s *Storage) Delete(jar string, subnet string) error {
	sql := `DELETE FROM %s WHERE sub_net = $1`
	sql = fmt.Sprintf(sql, jar)
	_, err := s.db.ExecContext(*s.ctx, sql, subnet)
	return err
}

func (s *Storage) Get(jar string, subnet string) (*common.IPSubnet, error) {
	IPSubnet := common.IPSubnet{}
	sql := `SELECT "sub_net","date_create" FROM %s WHERE sub_net = $1`
	sql = fmt.Sprintf(sql, jar)
	err := s.db.GetContext(*s.ctx, &IPSubnet, sql, subnet)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return &IPSubnet, common.ErrIPSubnetNotFound
	}
	return &IPSubnet, err
}

func (s *Storage) List(jar string) ([]*common.IPSubnet, error) {
	IPSubnet := make([]*common.IPSubnet, 0)
	sql := `SELECT "sub_net","date_create" FROM %s`
	sql = fmt.Sprintf(sql, jar)
	err := s.db.SelectContext(*s.ctx, &IPSubnet, sql)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return IPSubnet, common.ErrIPSubnetNotFound
	}
	return IPSubnet, err
}

func (s *Storage) Clear(jar string) error {
	sql := `TRUNCATE TABLE %s`
	sql = fmt.Sprintf(sql, jar)
	_, err := s.db.ExecContext(*s.ctx, sql)
	return err
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

func (s *Storage) IsOverlapping(jar string, subnet *common.IPSubnet) (bool, error) {
	count := 0
	sql := `SELECT count(*) FROM %s WHERE $1::cidr << "sub_net"::cidr or "sub_net"::cidr << $1::cidr 
                            or "sub_net"::cidr = $1::cidr`
	sql = fmt.Sprintf(sql, jar)
	err := s.db.GetContext(*s.ctx, &count, sql, subnet.Subnet)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Storage) InSubNet(jar string, ip string) (bool, error) {
	count := 0
	sql := `SELECT count(*) FROM %s WHERE $1::inet << "sub_net"::cidr`
	sql = fmt.Sprintf(sql, jar)
	err := s.db.GetContext(*s.ctx, &count, sql, ip)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Storage) Load(_ string) (bool, error) {
	return true, nil
}
