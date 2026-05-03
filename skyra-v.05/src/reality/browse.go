package reality

import (
	"os"
	"strings"

	firecrawl "github.com/mendableai/firecrawl-go"
)

type Browse struct {
	id string
}

func (b *Browse) ID() string { return b.id }

func (b *Browse) Create(r *Relation) Reality {
	return &Browse{id: "browse"}
}

func (b *Browse) Realize(r *Relation) string {
	url := strings.TrimSpace(r.Impulse)
	if url == "" {
		return "no url provided"
	}

	app, err := firecrawl.NewFirecrawlApp(os.Getenv("FIRECRAWL_API_KEY"), "")
	if err != nil {
		if r.Log != nil {
			r.Log("[browse]: init error:", err)
		}
		return "browse unavailable"
	}

	doc, err := app.ScrapeURL(url, nil)
	if err != nil {
		if r.Log != nil {
			r.Log("[browse]: scrape error:", err)
		}
		return "failed to fetch: " + url
	}

	content := doc.Markdown
	if len(content) > 4000 {
		content = content[:4000] + "\n\n[truncated]"
	}

	if r.Log != nil {
		r.Log("[browse]: fetched", url, "| length:", len(content))
	}
	return content
}
