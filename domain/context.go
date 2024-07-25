package domain

import "context"

type IContext interface {
	context.Context
}

type ContextKey string
