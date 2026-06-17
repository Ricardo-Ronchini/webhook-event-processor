package contexts

import (
	"context"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
)

type Context struct {
	app           *App
	helper        *Helper
	Host          string
	APIVersion    string
	SystemContext context.Context
}

func NewContext() *Context {
	return &Context{
		app:           NewApp(),
		helper:        &Helper{},
		Host:          common.GetEnv("HOST", "localhost:8080"),
		APIVersion:    common.GetEnv("API_VERSION", "V.0"),
		SystemContext: context.Background(),
	}
}

func (c *Context) App() *App {
	return c.app
}

func (c *Context) Helper() *Helper {
	return c.helper
}
