//go:build !windows

package activewindow

import (
	"context"

	"evolord-client/cmd/agent/runtime"
)

func Start(_ context.Context, _ *runtime.Env) error {
	return nil
}
