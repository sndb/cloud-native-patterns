package patterns

import "context"

type Circuit func(context.Context) (string, error)
type Effector func(context.Context) (string, error)
type SlowFunction func(string) (string, error)
type WithContext func(context.Context, string) (string, error)
