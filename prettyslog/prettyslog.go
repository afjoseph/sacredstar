package prettyslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type PrettyJSONHandler struct {
	out  io.Writer
	opts *slog.HandlerOptions
	mu   sync.Mutex
}

func NewPrettyJSONHandler(
	w io.Writer,
	opts *slog.HandlerOptions,
) *PrettyJSONHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyJSONHandler{
		out:  w,
		opts: opts,
	}
}

func (h *PrettyJSONHandler) Enabled(
	ctx context.Context,
	level slog.Level,
) bool {
	if h.opts.Level == nil {
		return true
	}
	return level >= h.opts.Level.Level()
}

func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Create a map to store all attributes
	attrs := make(map[string]interface{})

	// Collect all attributes
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	// Use spew to dump the attributes
	fmt.Fprintf(h.out, "%s [%s] %s: %s",
		r.Time.Format("2006-01-02 15:04:05.000000"),
		r.Level.String(),
		r.Message,
		spew.Sdump(attrs),
	)
	return nil
}

func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler {
	return h
}
