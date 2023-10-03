package main

import (
	"context"
)

type Flyio struct {}

func (m *Flyio) MyFunction(ctx context.Context, stringArg string) (*Container, error) {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg}).Sync(ctx)
}
