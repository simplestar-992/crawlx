# CrawlX | Fast Web Crawler

<p align="center">
  <img src="https://img.shields.io/badge/Scraping-Web%20Crawler-F39C12?style=for-the-badge" alt=""/>
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt=""/>
</p>

---

### High-performance web crawling & scraping

CrawlX extracts links, pages, and structured data from websites at blazing speed.

```bash
crawlx https://example.com -depth 3 -workers 10
```

---

## Features

- 🚀 **Concurrent crawling** - Multiple workers, fast results
- 🎯 **Depth control** - Limit how deep to crawl
- 🔗 **Link extraction** - Find all links on a page
- 📄 **Content extraction** - Pull title, text, images
- 📊 **Sitemap generation** - Create sitemaps automatically
- 📝 **Export** - JSON, CSV, or plain text

---

## Installation

```bash
git clone https://github.com/simplestar-992/crawlx.git
cd crawlx
go build -o crawlx -ldflags="-s -w"
```

---

## Usage

```bash
# Basic crawl
crawlx https://example.com

# Custom depth and workers
crawlx https://example.com -depth 5 -workers 20

# Extract specific content
crawlx https://blog.example.com -extract "article-title,article-body"

# Output formats
crawlx https://example.com -format json -o results.json
```

---

## Examples

```bash
# Site audit
crawlx https://example.com -depth 2 | grep -i "login\|admin"

# Find all external links
crawlx https://example.com -extract links | grep "^http"

# Generate sitemap
crawlx https://example.com -sitemap -o sitemap.xml
```

---

MIT © 2024 [simplestar-992](https://github.com/simplestar-992)
