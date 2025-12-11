/*
Syndication (feed) generator library for golang.

Installing

	go get github.com/gorilla/feeds

Feeds provides a simple, generic Feed interface with a generic Item object as well as RSS, Atom, OPML and JSON Feed specific RssFeed, AtomFeed, OpmlFeed and JSONFeed objects which allow access to all of each spec's defined elements.

# Examples

Create a Feed and some Items in that feed using the generic interfaces:

	import (
		"time"
		. "github.com/gorilla/feeds"
	)

	now = time.Now()

	feed := &Feed{
		Title:       "jmoiron.net blog",
		Link:        &Link{Href: "http://jmoiron.net/blog"},
		Description: "discussion about tech, footie, photos",
		Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		Created:     now,
		Copyright:   "This work is copyright © Benjamin Button",
	}

	feed.Items = []*Item{
		&Item{
			Title:       "Limiting Concurrency in Go",
			Link:        &Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Description: "A discussion on controlled parallelism in golang",
			Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     now,
		},
		&Item{
			Title:       "Logic-less Template Redux",
			Link:        &Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
			Description: "More thoughts on logicless templates",
			Created:     now,
		},
		&Item{
			Title:       "Idiomatic Code Reuse in Go",
			Link:        &Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
			Description: "How to use interfaces <em>effectively</em>",
			Created:     now,
		},
	}

From here, you can output Atom, RSS, OPML, or JSON Feed versions of this feed easily

	atom, err := feed.ToAtom()
	rss, err := feed.ToRss()
	opml, err := feed.ToOpml()
	json, err := feed.ToJSON()

You can also get access to the underlying objects that feeds uses to export its XML

	atomFeed := (&Atom{Feed: feed}).AtomFeed()
	rssFeed := (&Rss{Feed: feed}).RssFeed()
	opmlFeed := (&Opml{Feed: feed}).OpmlFeed()
	jsonFeed := (&JSON{Feed: feed}).JSONFeed()

From here, you can modify or add each syndication's specific fields before outputting

	atomFeed.Subtitle = "plays the blues"
	atom, err := ToXML(atomFeed)
	rssFeed.Generator = "gorilla/feeds v1.0 (github.com/gorilla/feeds)"
	rss, err := ToXML(rssFeed)
	opmlFeed.Head.OwnerId = "abcd1234"
	atom, err := ToXML(atomFeed)
	jsonFeed.NextUrl = "https://www.example.com/feed.json?page=2"
	json, err := jsonFeed.ToJSON()
*/
package feeds
