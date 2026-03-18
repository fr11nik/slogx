// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// The original code is from:
// 1.blog https://avito.tech/old/tpost/yfhycud6h1-logirovanie-kak-v-avito-goslog
// 2.code https://github.com/avito-tech/avitotech-presentations/commit/cf5ff7ea041dcdd3846634239d9ac27d5c80a86a
package slogx

import (
	"io"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Option func(*MultiHandler)

// WithOtelSlogOption adds OpenTelemetry slog bridge handler.
func WithOtelSlogOption(serviceName string, opts ...otelslog.Option) Option {
	return func(h *MultiHandler) {
		tt := otelslog.NewHandler(serviceName, opts...)
		h.AddHandler(tt)
	}
}

// WithJSONHandler adds slog.JSONHandler as the primary handler.
func WithJSONHandler(w io.Writer, opts *slog.HandlerOptions) Option {
	return func(h *MultiHandler) {
		if opts == nil {
			opts = &slog.HandlerOptions{Level: slog.LevelDebug}
		}
		handler := slog.Handler(slog.NewJSONHandler(w, opts))
		handler = NewHandlerMiddleware(handler)
		h.AddHandler(handler)
	}
}

// WithTextHandler adds slog.TextHandler as the primary handler.
func WithTextHandler(w io.Writer, opts *slog.HandlerOptions) Option {
	return func(h *MultiHandler) {
		if opts == nil {
			opts = &slog.HandlerOptions{Level: slog.LevelDebug}
		}
		handler := slog.Handler(slog.NewTextHandler(w, opts))
		handler = NewHandlerMiddleware(handler)
		h.AddHandler(handler)
	}
}

// WithHandler adds any custom slog.Handler.
// The handler is wrapped with HandlerMiddleware for context extraction.
func WithHandler(handler slog.Handler) Option {
	return func(h *MultiHandler) {
		h.AddHandler(NewHandlerMiddleware(handler))
	}
}

// WithRawHandler adds any custom slog.Handler WITHOUT wrapping in HandlerMiddleware.
// Use this when you handle context extraction yourself.
func WithRawHandler(handler slog.Handler) Option {
	return func(h *MultiHandler) {
		h.AddHandler(handler)
	}
}
