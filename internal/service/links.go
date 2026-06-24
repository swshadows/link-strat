package links

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/go-chi/render"
)

type LinkResponse struct {
	Url     string `json:"url"`
	Status  string `json:"status"`
	IsAlive bool   `json:"is_alive"`
}

type LinkRequest struct {
	Urls []string `json:"urls"`
}

type LinkHandler struct {
	LinkRe *regexp.Regexp
}

const workers = 20

func NewLinkHandler() *LinkHandler {
	return &LinkHandler{
		LinkRe: regexp.MustCompile(`(https?[^\s()<>\"']+)`),
	}
}

func worker(ctx context.Context, jobs <-chan string, results chan<- LinkResponse, client *http.Client) {

	for link := range jobs {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, link, nil)
		if err != nil {
			results <- LinkResponse{Url: link, Status: err.Error(), IsAlive: false}
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			results <- LinkResponse{Url: link, Status: err.Error(), IsAlive: false}
			continue
		}
		resp.Body.Close()

		results <- LinkResponse{Url: link, Status: resp.Status, IsAlive: true}
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

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	jobs := make(chan string)
	results := make(chan LinkResponse)
	for range workers {
		go worker(r.Context(), jobs, results, &client)
	}
	go func() {
		defer close(jobs)
		for _, link := range filteredLinks {
			jobs <- link
		}
	}()

	linksResponse := make([]LinkResponse, 0, len(filteredLinks))
	for range len(filteredLinks) {
		linksResponse = append(linksResponse, <-results)
	}

	close(results)
	started := time.Now()
	render.JSON(w, r, map[string]any{"took_ms": time.Since(started).Milliseconds(), "ok": true, "links": linksResponse})
}

func (h *LinkHandler) validateLink(lnk string) bool {
	parsedUrl, err := url.Parse(lnk)

	return err == nil && (parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https")
}
