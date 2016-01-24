package allrecipes

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gophergala2016/recipe/pkg/repo"
	"golang.org/x/net/context"
	"strings"
)

const (
	mainpage = "http://allrecipes.com"
)

type job struct {
	url     string
	err     error
	doc     *goquery.Document
	results []*repo.RecipeLink
}

func index(ctx context.Context, lcache repo.LocalCache) {
	fetchQueue := make(chan *job)
	parseQueue := make(chan *job)
	cacheQueue := make(chan *job)

	queueCtx, cancelQueueing := context.WithCancel(ctx)

	// TODO support multiple fetch goroutines
	// TODO Control closing of goroutines with contexts not channel closing

	go queue(queueCtx, fetchQueue)
	go fetch(ctx, fetchQueue, parseQueue)
	go parse(ctx, parseQueue, cacheQueue)
	cache(ctx, cacheQueue, lcache, cancelQueueing)
	return
}

func queue(ctx context.Context, out chan<- *job) {
	defer fmt.Println("DEBUG: queue returning")
	page := 1
	for {
		j := &job{
			url: fmt.Sprintf("%s/recipes/?sort=Newest&page=%d", mainpage, page),
		}
		select {
		case <-ctx.Done():
			close(out)
			return
		case out <- j:
			fmt.Printf("Queued: %s\n", j.url)
			page++
		}
	}
}

func fetch(ctx context.Context, in <-chan *job, out chan<- *job) {
	defer fmt.Println("DEBUG: fetch returning")
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-in:
			if !ok {
				close(out)
				return
			}
			j.doc, j.err = goquery.NewDocument(j.url)
			if j.err != nil {
				// TODO Do something with the error
				continue
			}
			out <- j
		}
	}
}

func parse(ctx context.Context, in <-chan *job, out chan<- *job) {
	defer fmt.Println("DEBUG: parse returning")
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-in:
			if !ok {
				close(out)
				return
			}
			doc := j.doc
			j.results = make([]*repo.RecipeLink, 0)
			articles := doc.Find("article")
			articles.Each(func(i int, article *goquery.Selection) {
				as := article.ChildrenFiltered("a[data-internal-referrer-link='recipe hub']")
				a := as.First()
				url, exists := a.Attr("href")
				if !exists {
					return
				}

				titleSelection := article.Find("ar-save-item")
				title, exists := titleSelection.Attr("data-name")
				if exists {
					title = strings.Trim(title, "\"")
				}

				descriptionSel := article.Find("div[class='rec-card__description']")
				description := descriptionSel.Text()

				j.results = append(j.results, &repo.RecipeLink{
					Title:       title,
					Description: description,
					URL:         url,
				})
			})
			out <- j
		}
	}
}

func cache(ctx context.Context, in <-chan *job, cache repo.LocalCache, cancelQueueing func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-in:
			if !ok {
				return
			}
			if len(j.results) == 0 {
				cancelQueueing()
			}
			hasNewResult := false
			for _, res := range j.results {
				if cache.Cached(res) {
					continue
				}
				cache.Add(res)
				hasNewResult = true
			}
			if !hasNewResult {
				cancelQueueing()
			}
		}
	}
}
