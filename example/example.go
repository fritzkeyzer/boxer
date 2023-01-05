package main

import "github.com/fritzkeyzer/boxer"

func main() {
	//b := boxer.NewWithBorders()
	//b.WrapLimit = 50
	//b.AddHeader("id", "column 1", "two", "three")
	//b.AddLine(1, "foo", "bar", "baz")
	//b.AddLine(2, "foo", "bar", "baz")
	//b.AddLine(3, "multi line text", "is\nhandled\ncorrectly")
	//b.AddLine(4, "long lines are wrapped", "to fit into box max width", "box width is specified through WrapLimit")
	//b.Print()
	//┌──┬────────────────┬────────────────┬────────────────┐
	//│id│column 1        │two             │three           │
	//├──┼────────────────┼────────────────┼────────────────┤
	//│1 │foo             │bar             │baz             │
	//│2 │foo             │bar             │baz             │
	//│3 │multi line text │is              │                │
	//│  │                │handled         │                │
	//│  │                │correctly       │                │
	//│4 │long lines are  │to fit into box │box width is    │
	//│  │wrapped         │max width       │specified       │
	//│  │                │                │through         │
	//│  │                │                │WrapLimit       │
	//└──┴────────────────┴────────────────┴────────────────┘

	//id   column 1          two               three
	//1    foo               bar               baz
	//2    foo               bar               baz
	//3    multi line text   is
	//                       handled
	//                       correctly
	//4    long lines are    to fit into box   box width is
	//     wrapped           max width         specified
	//                                         through
	//                                         WrapLimit

	b := boxer.NewWithBorders()
	b.WrapLimit = 50
	b.AddHeader("id", "name", "value")
	b.AddLine(1, "foo", "bar")
	b.AddLine(2, "example", "multi\nline")
	b.Print()

}
