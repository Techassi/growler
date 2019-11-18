package growler

import (
	"golang.org/x/net/html"
)

type CollectorHTMLNode struct {
	Name         string
	Collector   *Collector
	attributes []html.Attribute
}

func (n *CollectorHTMLNode) Attr(attr string) string {
	for _, a := range n.attributes {
		if a.Key == attr {
			return a.Val
		}
	}

	return ""
}
