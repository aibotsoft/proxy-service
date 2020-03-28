package migration

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"path/filepath"
)

const (
	versionTable = "public.schema_version"
	pathFromRoot = "migrations"
)

// Up migrate db all the way up
func Up(conn *pgx.Conn) error {
	ctx := context.Background()
	m, err := prepareMigrations(ctx, conn)
	if err != nil {
		return err
	}
	return m.Migrate(ctx)
}

func UpTo(conn *pgx.Conn, targetVersion int) error {
	ctx := context.Background()
	m, err := prepareMigrations(ctx, conn)
	if err != nil {
		return err
	}
	return m.MigrateTo(ctx, int32(targetVersion))

}

func prepareMigrations(ctx context.Context, conn *pgx.Conn) (*migrate.Migrator, error) {
	migrationsPath := filepath.Join(config.RootDir(), pathFromRoot)
	m, err := migrate.NewMigrator(ctx, conn, versionTable)
	if err != nil {
		return m, err
	}
	err = m.LoadMigrations(migrationsPath)
	if err != nil {
		return m, err
	}
	return m, nil
}

//if err != nil {
//return err
//}
//v, err := m.GetCurrentVersion(ctx)
//if err != nil {
//return err
//}
//if v != targetVersion {
//return errors.New("targetVersion != actual version")
//}
//return nil
