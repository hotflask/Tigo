package TigoWeb

import (
	"errors"
	"net/http"
)

// Middleware http中间件
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// chainMiddleware 是http中间件生成器
func chainMiddleware(mw ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

// InternalServerErrorMiddleware 用来处理控制层出现的异常的中间件
func InternalServerErrorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequestProcessTimeMiddleware 用来记录接口处理请求所用时间的中间件
func RequestProcessTimeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStart := time.Now().Nanosecond() / 1e6
		next.ServeHTTP(w, r)
		requestEnd := time.Now().Nanosecond() / 1e6
		logger.Trace.Printf("%s %s %dms", r.Method, r.RequestURI, requestEnd-requestStart)
	})
}
