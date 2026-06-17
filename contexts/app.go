package contexts

import (
	"os"

	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/cache"
	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/db"
	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/redpanda"
)

type App struct {
	db     *db.Database
	cache  cache.CacheClient
	logs   *Logs
	broker *redpanda.Redpanda
}

func NewApp() *App {
	var cacheClient cache.CacheClient
	if addr := os.Getenv("DRAGONFLY_ADDR"); addr != "" {
		cacheClient = cache.NewClient(addr)
	}

	return &App{
		db:     db.NewDatabase(),
		cache:  cacheClient,
		logs:   NewLogs(),
		broker: redpanda.NewRedpanda(),
	}
}

func (a *App) Database() *db.Database       { return a.db }
func (a *App) Cache() cache.CacheClient     { return a.cache }
func (a *App) Logs() *Logs                  { return a.logs }
func (a *App) Redpanda() *redpanda.Redpanda { return a.broker }
