package route

import (
	"fmt"
	"net/http"
	"time"
)

func logCall(requestHandler http.Handler, routeName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := newResponseWriterWrapper(w)
		start := time.Now()
		requestHandler.ServeHTTP(res, r)
		duration := time.Since(start)
		durationMs := duration.Seconds() * float64(time.Second/time.Millisecond)
		fmt.Printf(
			//"\t%s\t%s\t%s\t%s\t%d\t%f",
			"time: %s, method: %s, path: %s, route: %s, statusCode: %d, duration: %f\n",
			time.Now(),
			r.Method,
			r.RequestURI,
			routeName,
			res.Status(),
			durationMs,
		)
	})
}
