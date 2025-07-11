// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// The original code is from:
// 1.blog https://avito.tech/old/tpost/yfhycud6h1-logirovanie-kak-v-avito-goslog
// 2.code https://github.com/avito-tech/avitotech-presentations/commit/cf5ff7ea041dcdd3846634239d9ac27d5c80a86a
package slogx

import (
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Option func(*MultiHandler)

func WithOtelSlogOption(serviceName string, opts ...otelslog.Option) Option {
	return func(h *MultiHandler) {
		tt := otelslog.NewHandler(serviceName, opts...)
		h.handlers = append(h.handlers, tt)
	}
}
