// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// The original code is from:
// 1.blog https://avito.tech/old/tpost/yfhycud6h1-logirovanie-kak-v-avito-goslog
// 2.code https://github.com/avito-tech/avitotech-presentations/commit/cf5ff7ea041dcdd3846634239d9ac27d5c80a86a
package slogx

import (
	"context"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

// -----------------------------------------------

// HandlerMiddlware for middleware
type HandlerMiddlware struct {
	next slog.Handler
}

// NewHandlerMiddleware create new middleware
func NewHandlerMiddleware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}

// Enabled reports whether the handler handles records at the given level. The handler ignores records whose level is lower. It is called early, before any arguments are processed, to save effort if the log event should be discarded. If called from a Logger method, the first argument is the context passed to that method, or context.Background() if nil was passed or the method does not take a context. The context is passed so Enabled can use its values to make a decision.
func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

// Handle is called for each log record and add it to the record.
func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.UserID != uuid.Nil {
			rec.Add("userID", c.UserID)
		}
		if c.Phone != "" {
			rec.Add("phone", c.Phone)
		}
		if c.TraceID != "" {
			rec.Add("traceID", c.TraceID)
		}
		if c.Scope != "" {
			rec.Add("env", c.Scope)
		}
		if c.Message != "" {
			rec.Add("message", c.Message)
		}
		if c.AppName != "" {
			rec.Add("app_name", c.AppName)
		}
		if c.Language != "" {
			rec.Add("language", c.Language)
		}
		if c.Type != "" {
			rec.Add("type", c.Type)
		}
	}
	return h.next.Handle(ctx, rec)
}

// WithAttrs adds attributes to the log record.
func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

// WithGroup adds a new group to the log record.
func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}

type logCtx struct {
	SpanID   string
	UserID   uuid.UUID
	Phone    string
	Scope    string
	Type     string
	AppName  string
	Message  string
	Language string
	TraceID  string
}

type keyType int

const key = keyType(0)

// WithLogUserID by execution add userID to log record
func WithLogUserID(ctx context.Context, userID uuid.UUID) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}

func WithLogTraceID(ctx context.Context, traceID string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.TraceID = traceID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{TraceID: traceID})
}

// WithLogLanguage by execution add language to log record
func WithLogLanguage(ctx context.Context, language string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Language = language
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Language: language})
}

// WithLogAppName by execution add appName to log record
func WithLogAppName(ctx context.Context, appName string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.AppName = appName
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{AppName: appName})
}

// WithLogPhone by execution add phone to log record
func WithLogPhone(ctx context.Context, phone string) context.Context {
	if len(phone) > 4 {
		phone = strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
	}
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Phone = phone
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Phone: phone})
}

// WithLogScope by execution add scope to log record
func WithLogScope(ctx context.Context, envScope string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Scope = envScope
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Scope: envScope})
}

// WithLogType by execution add type to log record
func WithLogType(ctx context.Context, typelog string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Type = typelog
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Type: typelog})
}

// WithLogMessage by execution add message to log record
func WithLogMessage(ctx context.Context, message string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Message = message
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Message: message})
}

func WithSpanID(ctx context.Context, spanID string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Scope = spanID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{SpanID: spanID})
}

// -----------------------------------------------

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

// return next error
func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

// WrapError wraps err into ctx like linked list
func WrapError(ctx context.Context, err error) error {
	c := logCtx{}
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

// ErrorCtx adds ctx to err
func ErrorCtx(ctx context.Context, err error) context.Context {
	if e, ok := err.(*errorWithLogCtx); ok { // в реальной жизни используйте error.As
		return context.WithValue(ctx, key, e.ctx)
	}
	return ctx
}

func CombineContexts(logctx context.Context, requestCtx context.Context) context.Context {
	// Получаем logCtx из логирующего контекста
	if lc, ok := logctx.Value(key).(logCtx); ok {
		// Проверяем, если в requestCtx есть traceID, если нет — добавляем
		if rc, ok := requestCtx.Value(key).(logCtx); ok {
			if lc.TraceID == "" && rc.TraceID != "" {
				lc.TraceID = rc.TraceID
			}
		}
		// Возвращаем комбинированный контекст
		return context.WithValue(requestCtx, key, lc)
	}

	// Если logCtx нет, просто возвращаем requestCtx
	return requestCtx
}

// -----------------------------------------------

// InitLogging init logger by JSON hanlder and store into io.Writer
func InitLogging(w io.Writer, opts ...Option) {
	if w == nil {
		log.Fatal("slog: writer is nil")
	}
	hm := &MultiHandler{}
	for _, opt := range opts {
		opt(hm)
	}

	// Создаём мульти-обработчик
	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	handler = NewHandlerMiddleware(handler)
	hm.handlers = append(hm.handlers, handler)
	slog.SetDefault(slog.New(hm))
}
