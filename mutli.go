package slogx

import (
	"context"
	"log/slog"
	"sync"
)

// MultiHandler объединяет несколько slog.Handler и вызывает их при логировании.
type MultiHandler struct {
	mu       sync.RWMutex
	handlers []slog.Handler
}

// NewMultiHandler создаёт новый мульти-обработчик.
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

// AddHandler добавляет хендлер в рантайме (thread-safe).
// Позволяет inject'ить хендлеры после создания логгера.
func (mh *MultiHandler) AddHandler(h slog.Handler) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	mh.handlers = append(mh.handlers, h)
}

// Enabled возвращает true, если хотя бы один из обработчиков включён для данного уровня.
func (mh *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	mh.mu.RLock()
	defer mh.mu.RUnlock()
	for _, h := range mh.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle вызывает Handle для каждого из обработчиков.
func (mh *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	mh.mu.RLock()
	defer mh.mu.RUnlock()
	for _, h := range mh.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// WithAttrs возвращает новый MultiHandler с добавленными атрибутами для каждого из обработчиков.
func (mh *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	mh.mu.RLock()
	defer mh.mu.RUnlock()
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

// WithGroup возвращает новый MultiHandler с добавленной группой для каждого из обработчиков.
func (mh *MultiHandler) WithGroup(name string) slog.Handler {
	mh.mu.RLock()
	defer mh.mu.RUnlock()
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}
