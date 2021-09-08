package patterns

import "context"

type Circuit func(context.Context) (string, error)
type Effector func(context.Context) (string, error)
