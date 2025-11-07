package http_util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"

	config "oracle.com/oracle/my-go-oracle-app/configs"
)

// HTTPUtil interface
type HTTPUtil interface {
	Send(ctx context.Context, request *http.Request) (*http.Response, error)
}

type httpUtil struct {
	httpClient *http.Client
	debug      bool
}

// NewHTTPUtil create http service
func NewHTTPUtil(cfg *config.Config, duration string) HTTPUtil {

	var timeout time.Duration
	value, err := time.ParseDuration(duration)
	if err != nil {
		timeout = cfg.HTTPTimeout //default
	}
	timeout = value

	httpClient := &http.Client{
		Timeout: time.Duration(timeout),
	}
	httpClient.Transport = &http.Transport{
		MaxIdleConns:        cfg.HTTPMaxIdleConnections,
		MaxIdleConnsPerHost: cfg.HTTPMaxIdleConnectionsPerHost,
		IdleConnTimeout:     cfg.HTTPIdleConnectionTimeout,
	}
	return &httpUtil{
		httpClient: httpClient,
		debug:      cfg.HTTPDebug,
	}
}

func (h *httpUtil) send(ctx context.Context, request *http.Request) (response *http.Response, err error) {
	return h.httpClient.Do(request.WithContext(ctx))
}

func (h *httpUtil) Send(ctx context.Context, request *http.Request) (response *http.Response, err error) {
	response, err = h.send(ctx, request)
	h.dump(request, response)
	return
}

func (h *httpUtil) dump(request *http.Request, response *http.Response) {
	if !h.debug {
		return
	}

	slog.Info(string(DumpHTTPRequest(request)))
	if response == nil {
		return
	}

	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		return
	}

	slog.Info(fmt.Sprintf("Response: %q", responseDump))
}

// DumpHTTPRequest dump request without header, just body
func DumpHTTPRequest(req *http.Request) []byte {
	if req.Body == nil {
		return nil
	}

	save, body, err := drainBody(req.Body)
	if err != nil {
		return nil
	}

	req.Body = body

	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	var b bytes.Buffer

	if req.Body != nil {
		var dest io.Writer = &b
		if chunked {
			dest = httputil.NewChunkedWriter(dest)
		}
		_, err = io.Copy(dest, req.Body)
		if chunked {
			dest.(io.Closer).Close()
			io.WriteString(&b, "")
		}
	}

	req.Body = save
	if err != nil {
		return nil
	}

	return b.Bytes()
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}

	if err = b.Close(); err != nil {
		return nil, b, err
	}

	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
