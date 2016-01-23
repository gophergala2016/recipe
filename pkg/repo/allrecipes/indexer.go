package allrecipes

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/context"
	"strings"
	"sync"
)

const (
	mainpage = "http://allrecipes.com"
)

type entry struct {
	title       string
	description string
	url         string
}

func (e *entry) Title() string {
	return e.title
}

func (e *entry) Description() string {
	return e.description
}

func (e *entry) URL() string {
	return e.url
}

func index(ctx context.Context) <-chan *entry {
	limit := 10
	queueCtx, queueCancel := context.WithCancel(ctx)
	f := newFetcher(ctx, limit)
	entries := make(chan *entry)
	go parseLinks(ctx, queueCancel, f.responses(), entries)
	go queueLinkFetching(queueCtx, f)
	return entries
}

func queueLinkFetching(ctx context.Context, f *fetcher) {
	defer f.close()
	var (
		url  string
		page = 1
	)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		url = fmt.Sprintf("%s/recipes/?sort=Newest&page=%d", mainpage, page)
		if err := f.fetch(ctx, url); err != nil {
			return
		}
		page++
	}
}

type response struct {
	doc *goquery.Document
	err error
}

func parseLinks(ctx context.Context, cancelQueueing func(), in <-chan response, out chan<- *entry) {
	defer close(out)
	parsed := 0
	debugLimit := 20
	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-in:
			if !ok {
				return
			}
			if resp.err != nil {
				continue
			}
			doc := resp.doc
			articles := doc.Find("article")
			docHasRecipeLinks := false
			articles.Each(func(i int, article *goquery.Selection) {
				as := article.ChildrenFiltered("a[data-internal-referrer-link='recipe hub']")
				a := as.First()
				url, exists := a.Attr("href")
				if !exists {
					return
				}
				docHasRecipeLinks = true

				titleSelection := article.Find("ar-save-item")
				title, exists := titleSelection.Attr("data-name")
				if exists {
					title = strings.Trim(title, "\"")
				}

				descriptionSel := article.Find("div[class='rec-card__description']")
				description := descriptionSel.Text()

				out <- &entry{
					title:       title,
					description: description,
					url:         url,
				}
			})
			parsed++
			if !docHasRecipeLinks || parsed == debugLimit {
				// Cancel queueing new fetch requests for links to recipes
				cancelQueueing()
			}
		}
	}
}

type fetcher struct {
	limit  int
	ctx    context.Context
	in     chan string
	out    chan response
	cancel func()
	wg     sync.WaitGroup
}

func newFetcher(ctx context.Context, limit int) *fetcher {
	fetchCtx, cancel := context.WithCancel(ctx)
	f := &fetcher{
		limit:  limit,
		ctx:    fetchCtx,
		cancel: cancel,
		in:     make(chan string),
		out:    make(chan response),
	}
	f.wg.Add(limit)
	for i := 0; i < limit; i++ {
		go f.run()
	}
	return f
}

func (f *fetcher) close() {
	f.cancel()
	f.wg.Wait()
	close(f.in)
	close(f.out)
}

func (f *fetcher) responses() <-chan response {
	return f.out
}

func (f *fetcher) fetch(ctx context.Context, url string) error {
	select {
	case f.in <- url:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (f *fetcher) run() {
	defer f.wg.Done()
	for {
		select {
		case req := <-f.in:
			fetchDoc(f.ctx, req, f.out)
		case <-f.ctx.Done():
			return
		}
	}
}

func fetchDoc(ctx context.Context, url string, out chan<- response) {
	doc, err := goquery.NewDocument(url)
	select {
	case out <- response{doc, err}:
		return
	case <-ctx.Done():
		return
	}
}
