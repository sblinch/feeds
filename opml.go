package feeds

// rss support
// validation done according to spec here:
//    http://cyber.law.harvard.edu/rss/rss.html

import (
	"encoding/xml"
	"time"
)

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type OpmlFeed struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    *OpmlHead
	Body    *OpmlBody
}

type OpmlHead struct {
	XMLName         xml.Name `xml:"head"`
	Title           string   `xml:"title,omitempty"` // required
	DateCreated     string   `xml:"dateCreated,omitempty"`
	DateModified    string   `xml:"dateModified,omitempty"`
	OwnerName       string   `xml:"ownerName,omitempty"`
	OwnerEmail      string   `xml:"ownerEmail,omitempty"`
	OwnerId         string   `xml:"ownerId,omitempty"`         // the http address of a web page that contains information that allows a human reader to communicate with the author of the document via email or other means. It also may be used to identify the author. No two authors have the same ownerId.
	Docs            string   `xml:"docs,omitempty"`            //  the http address of documentation for the format used in the OPML file. It's probably a pointer to this page for people who might stumble across the file on a web server 25 years from now and wonder what it is.
	ExpansionState  string   `xml:"expansionState,omitempty"`  //  a comma-separated list of line numbers that are expanded. The line numbers in the list tell you which headlines to expand. The order is important. For each element in the list, X, starting at the first summit, navigate flatdown X times and expand. Repeat for each element in the list.
	VertScrollState int      `xml:"vertScrollState,omitempty"` // a number, saying which line of the outline is displayed on the top line of the window. This number is calculated with the expansion state already applied.
	WindowTop       int      `xml:"windowTop,omitempty"`       // a number, the pixel location of the top edge of the window.
	WindowLeft      int      `xml:"windowLeft,omitempty"`      // a number, the pixel location of the left edge of the window.
	WindowBottom    int      `xml:"windowBottom,omitempty"`    // a number, the pixel location of the bottom edge of the window.
	WindowRight     int      `xml:"windowRight,omitempty"`     // a number, the pixel location of the right edge of the window.
}

type OpmlBody struct {
	XMLName  xml.Name       `xml:"body"`
	Outlines []*OpmlOutline `xml:"outline"`
}

type OpmlOutline struct {
	XMLName      xml.Name       `xml:"outline"`
	Text         string         `xml:"text,attr"`                   // required
	Type         string         `xml:"type,attr,omitempty"`         // a string, it says how the other attributes of the <outline> are interpreted.#
	IsComment    string         `xml:"isComment,attr,omitempty"`    // a string, either "true" or "false", indicating whether the outline is commented or not. By convention if an outline is commented, all subordinate outlines are considered to also be commented. If it's not present, the value is false.
	IsBreakpoint string         `xml:"isBreakpoint,attr,omitempty"` // a string, either "true" or "false", indicating whether a breakpoint is set on this outline. This attribute is mainly necessary for outlines used to edit scripts. If it's not present, the value is false.#
	Created      string         `xml:"created,attr,omitempty"`
	Category     string         `xml:"category,attr,omitempty"` //  a string of comma-separated slash-delimited category strings, in the format defined by the RSS 2.0 category element. To represent a "tag," the category string should contain no slashes. Examples: 1. category="/Boston/Weather". 2. category="/Harvard/Berkman,/Politics".
	Outlines     []*OpmlOutline `xml:"outline,omitempty"`
	OpmlInclusion
	OpmlSubscriptionList
}

type OpmlSubscriptionList struct {
	// type=rss
	// the text attribute should initially be the top-level title element in the feed being pointed to, however since it is user-editable, processors should not depend on it always containing the title of the feed
	Title       string `xml:"title,attr,omitempty"`       // the top-level title element from the feed; probably the same as text, it should not be omitted
	Description string `xml:"description,attr,omitempty"` // the top-level description element from the feed. htmlUrl is the top-level link element. language is the value of the top-level language element. title is probably the same as text, it should not be omitted. title contains the top-level title element from the feed
	XmlUrl      string `xml:"xmlUrl,attr,omitempty"`      // the http address of the feed
	HtmlUrl     string `xml:"htmlUrl,attr,omitempty"`     // the top-level link element
	Language    string `xml:"language,attr,omitempty"`    // the value of the top-level language element
	Version     string `xml:"version,attr,omitempty"`     //  varies depending on the version of RSS that's being supplied. It was invented at a time when we thought there might be some processors that only handled certain versions, but that hasn't turned out to be a major issue. The values it can have are: RSS1 for RSS 1.0; RSS for 0.91, 0.92 or 2.0; scriptingNews for scriptingNews format
}

type OpmlInclusion struct {
	// type=link or type=include
	// the text element is, as usual, what's displayed in the outliner; it's also what is displayed in an HTML rendering
	Url string `xml:"url,attr,omitempty"`
}

type Opml struct {
	*Feed
}

// create a new Subscription List OpmlOutline with a generic Item struct's data
func newOpmlSubscriptionList(i *Item) *OpmlOutline {
	item := &OpmlOutline{
		Text:    i.Title,
		Type:    "rss",
		Created: anyTimeFormat(time.RFC822, i.Created, i.Updated),
		OpmlSubscriptionList: OpmlSubscriptionList{
			Title:       i.Title,
			Description: i.Description,
		},
	}
	if i.Source != nil {
		item.XmlUrl = i.Source.Href
		if i.Link != nil {
			item.HtmlUrl = i.Link.Href
		}
	} else if i.Link != nil {
		item.XmlUrl = i.Link.Href
	}
	return item
}

// create a new Inclusion OpmlOutline with a generic Item struct's data
func newOpmlInclusion(i *Item) *OpmlOutline {
	item := &OpmlOutline{
		Text:    i.Title,
		Type:    "link",
		Created: anyTimeFormat(time.RFC822, i.Created, i.Updated),
	}

	if i.Link != nil {
		item.Url = i.Link.Href
	}
	return item
}

// create a new RssFeed with a generic Feed struct's data
func (o *Opml) OpmlFeed() *OpmlFeed {
	feed := &OpmlFeed{
		Version: "2.0",
		Head: &OpmlHead{
			Title:        o.Title,
			DateCreated:  anyTimeFormat(time.RFC822, o.Created, o.Updated),
			DateModified: anyTimeFormat(time.RFC822, o.Updated),
		},
		Body: &OpmlBody{},
	}

	if o.Author != nil {
		feed.Head.OwnerName = o.Author.Name
		feed.Head.OwnerEmail = o.Author.Email
		feed.Head.OwnerId = o.Author.Email
	}

	for _, i := range o.Items {
		feed.Body.Outlines = append(feed.Body.Outlines, newOpmlInclusion(i))
	}
	return feed
}

// FeedXml returns an XML-Ready object for an Rss object
func (o *Opml) FeedXml() interface{} {
	// only generate version 2.0 feeds for now
	return o.OpmlFeed().FeedXml()

}

// FeedXml returns an XML-ready object for an RssFeed object
func (r *OpmlFeed) FeedXml() interface{} {
	return r
}
