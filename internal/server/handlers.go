package server

import (
	"io"
	"maps"
	"net/http"
)

// simple passthrough for non-cacheable requests
func (p *ProxyServer) forwardToOrigin(
	w http.ResponseWriter, r *http.Request, targetURL string,
) {
	resp := p.fetchFromOrigin(r, w, targetURL)
	if resp == nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	sendFinalResponse(w, resp.Header, "MISS", resp.StatusCode, body)
	p.logger.Infof("'%s %s' is served through origin (non-cacheable)", r.Method, targetURL)
}

func (p *ProxyServer) fetchFromOriginAndCacheThenServe(
	w http.ResponseWriter, r *http.Request, targetURL string,
) {
	// get from origin
	resp := p.fetchFromOrigin(r, w, targetURL)
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	// first cache the response (if cacheable)
	if isCacheable, cacheTTl := isResponseCacheable(resp); isCacheable {
		if err := p.cacheService.Save(p.ctx, resp, body, targetURL, cacheTTl); err != nil {
			p.logger.Errorf(
				"failed to cache response for '%s %s'. error message: %s", r.Method, targetURL, err,
			)
		} else {
			p.logger.Infof("response for '%s %s' cached successfully!", r.Method, targetURL)
		}
	}

	// then send response to client
	sendFinalResponse(w, resp.Header, "MISS", resp.StatusCode, body)
	p.logger.Infof("'%s %s' is served through origin.", r.Method, targetURL)
}

func (p *ProxyServer) fetchFromOrigin(
	r *http.Request, w http.ResponseWriter, targetURL string,
) *http.Response {
	// prepare outgoing request (Origin Server Connection Manager)
	outReq, _ := http.NewRequest(r.Method, targetURL, r.Body)
	maps.Copy(outReq.Header, r.Header)

	response, err := p.httpClient.Do(outReq)
	if err != nil {
		http.Error(w, "Origin Unreachable", http.StatusBadGateway)
		p.logger.Infof("origin unreachable for '%s %s'. Error message: %s", r.Method, targetURL, err)
		return nil
	}

	return response
}

func sendFinalResponse(w http.ResponseWriter, headers http.Header, XCache string, statusCode int, body []byte) {
	maps.Copy(w.Header(), headers)
	w.Header().Set("X-Cache", XCache)
	w.WriteHeader(statusCode)
	w.Write(body)
}
