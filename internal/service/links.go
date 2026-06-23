package links

import (
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/go-chi/render"
)

type Link struct {
	Url    string `json:"url"`
	Status string `json:"status"`
}

type LinkRequest struct {
	Urls []string `json:"urls"`
}

type LinkHandler struct {
	LinkRe *regexp.Regexp
}

func NewLinkHandler() *LinkHandler {
	return &LinkHandler{
		LinkRe: regexp.MustCompile(`(https?[^\s()<>\"']+)`),
	}
}

func (h *LinkHandler) CheckLinks(w http.ResponseWriter, r *http.Request) {
	var linkRequest LinkRequest

	if err := render.DecodeJSON(r.Body, &linkRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]any{"ok": false, "error": "missing links in request body"})
		return
	}

	filteredLinks := make([]string, 0)

	for _, link := range linkRequest.Urls {
		if h.validateLink(link) {
			filteredLinks = append(filteredLinks, link)
		}
	}

	if len(filteredLinks) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]any{"ok": false, "error": "no valid links in request body"})
		return
	}

	var linksResponse []Link

	httpClient := &http.Client{Timeout: 10 * time.Second}
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	started := time.Now()
	for _, link := range filteredLinks {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			resp, err := httpClient.Get(link)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			mu.Lock()
			linksResponse = append(linksResponse, Link{Url: link, Status: resp.Status})
			mu.Unlock()
		}(link)
	}
	wg.Wait()

	render.JSON(w, r, map[string]any{"took_ms": time.Since(started).Milliseconds(), "ok": true, "links": linksResponse})
}

func (h *LinkHandler) validateLink(lnk string) bool {
	parsedUrl, err := url.Parse(lnk)

	return err == nil && (parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https")
}
