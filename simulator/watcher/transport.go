package watcher

import "net/http"

type dockerRewriteTransport struct {
	rt http.RoundTripper
}

func (t *dockerRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// リダイレクト先が localhost:810x だったら Docker 内部ホスト名に書き換える
	switch req.URL.Host {
	case "localhost:8101":
		req.URL.Host = "delivery-node-1:80"
	case "localhost:8102":
		req.URL.Host = "delivery-node-2:80"
	case "localhost:8103":
		req.URL.Host = "delivery-node-3:80"
	}
	return t.rt.RoundTrip(req)
}
