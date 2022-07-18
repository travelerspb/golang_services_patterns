package src

import "context"

type Circuit func(context.Context) (string, error)
