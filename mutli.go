// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// The original code is from:
// 1.blog https://avito.tech/old/tpost/yfhycud6h1-logirovanie-kak-v-avito-goslog
// 2.code https://github.com/avito-tech/avitotech-presentations/commit/cf5ff7ea041dcdd3846634239d9ac27d5c80a86a
package slogx

import (
	"context"
	"log/slog"
)

// MultiHandler объединяет несколько slog.Handler и вызывает их при логировании.
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler создаёт новый мульти-обработчик.
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

// Enabled возвращает true, если хотя бы один из обработчиков включён для данного уровня.
func (mh *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range mh.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle вызывает Handle для каждого из обработчиков.
func (mh *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range mh.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// WithAttrs возвращает новый MultiHandler с добавленными атрибутами для каждого из обработчиков.
func (mh *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

// WithGroup возвращает новый MultiHandler с добавленной группой для каждого из обработчиков.
func (mh *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}
