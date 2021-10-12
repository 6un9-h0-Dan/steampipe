package db_common

import (
	"context"
	"database/sql"

	"github.com/turbot/steampipe/steampipeconfig"

	"github.com/turbot/steampipe/query/queryresult"
	"github.com/turbot/steampipe/schema"
)

type EnsureSessionStateCallback = func(context.Context, *sql.Conn) error

type Client interface {
	Close() error
	LoadSchema()

	SchemaMetadata() *schema.Metadata
	ConnectionMap() *steampipeconfig.ConnectionDataMap

	GetCurrentSearchPath() ([]string, error)
	SetSessionSearchPath(...string) error

	ExecuteSync(ctx context.Context, query string, disableSpinner bool) (*queryresult.SyncQueryResult, error)
	Execute(ctx context.Context, query string, disableSpinner bool) (res *queryresult.Result, err error)

	CacheOn() error
	CacheOff() error
	CacheClear() error

	SetEnsureSessionDataFunc(EnsureSessionStateCallback)
	RefreshSessions(ctx context.Context) error
	// remote client will have empty implementation

	RefreshConnectionAndSearchPaths() *RefreshConnectionResult
}
