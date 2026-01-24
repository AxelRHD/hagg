package shared

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MainContainer(content ...g.Node) g.Node {
	return Main(
		Class("container pt-0"),
		g.Group(content),
	)
}
