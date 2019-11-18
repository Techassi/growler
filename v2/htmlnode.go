package growler

type CollectorHTMLNode struct {
	Name         string
	Collector   *Collector
	attributes []html.Attribute
}

func (n *CollectorHTMLNode) Attr(attr string) string {
	for _, a := range n.attributes {
		if a.Key == k {
			return a.Val
		}
	}
	return ""
}
