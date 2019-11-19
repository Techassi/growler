# growler

growler is a web crawler written in Go and tuned to work in a multi-threaded enviroment.

This project is a personal one and was intended to get some experience in multi-threaded programming in Go.

## Usage

Start using growler by importing it into your project by simply adding this to your imports:

```golang
import (
	"github.com/Techassi/growler"
)
```

### Creating a collector

The collector struct is created like this:

```golang
c := growler.NewCollector()
```

The returend collector is used to customize and manage the crawling process.

### Using the collector

To start looking for URLs you can use the `OnHTML` hook, which gets triggered at a HTML response.

```golang
c.OnHTML("a[href]", func (n growler.CollectorHTMLNode) {
	// get the href attribute to crawl this as the next URL
	link := n.Attr("href")

	// Print out the href attribute
	fmt.Println(link)

	// Call the Visit function to visit the found URL
	n.Collector.Visit(link)
})
```

`OnHTML` expects two parameters, a HTML selector like `a[href]` and a callback function which gets executed for each element found by the provided selector. This callback function needs to be defined as `func (n growler.CollectorHTMLNode)`

### Customizing the collector
#### Delay

To set a delay (recommended) use the `SetDelay(d int)` function:

```golang
c.SetDelay(2)
```

The only parameter to be provided is an integer to specify the amount of seconds the delay should last.

#### MaxDepth

To set a max depth of URLs to be crawled use the `SetMaxDepth(d int)` function:

```golang
c.SetMaxDepth(2)
```

The only parameter to be provided is an integer to specify the max depth. Example: example.com/foo/bar is valid, example.com/foo/bar/baz is not and will NOT be crawled.

#### Set the seed / start URL

To set the seed / start URL use the `Visit(url string)` function:

```golang
c.Visit("https://example.com/")
```

### Complete example

```golang
c := growler.NewCollector()

c.OnHTML("a[href]", func (n growler.CollectorHTMLNode) {
	link := n.Attr("href")

	fmt.Println(link)

	n.Collector.Visit(link)
})

c.SetDelay(2)
c.SetMaxDepth(2)

c.Visit("https://example.com/")
c.Wait()
```
