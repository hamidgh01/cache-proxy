package cache

import "net/http"

var hopByHopHeaders = map[string]string{
	"Connection":          "",
	"Keep-Alive":          "",
	"Proxy-Authenticate":  "",
	"Proxy-Authorization": "",
	"TE":                  "",
	"Trailer":             "",
	"Transfer-Encoding":   "",
	"Upgrade":             "",
}

func filterResponseHeaders(h http.Header) http.Header {
	cleanedHeaders := make(http.Header)

	for key, values := range h {
		// skip hop-by-hop headers
		if _, ok := hopByHopHeaders[key]; ok {
			continue
		}

		for _, v := range values {
			cleanedHeaders.Add(key, v)
		}
	}

	return cleanedHeaders
}
