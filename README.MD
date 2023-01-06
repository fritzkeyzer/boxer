# Boxer
Is a simple table writer utility package that supports multiline text input 
and will automatically wrap long lines.

Example
```go
b := boxer.NewWithBorders()
b.WrapLimit = 40
b.AddHeader("id", "name", "value")
b.AddLine(1, "foo", "bar")
b.AddLine(2, "example", "multi\nline\ninput")
b.AddLine("columns can contain anything", "for example", "a very very long line that will be wrapped to maintain the overall box width")
b.Print()


//┌───────────────┬───────────┬──────────────┐
//│id             │name       │value         │
//├───────────────┼───────────┼──────────────┤
//│1              │foo        │bar           │
//│2              │example    │multi         │
//│               │           │line          │
//│               │           │input         │
//│columns can    │for example│a very very   │
//│contain        │           │long line     │
//│anything       │           │that will be  │
//│               │           │wrapped to    │
//│               │           │maintain the  │
//│               │           │overall box   │
//│               │           │width         │
//└───────────────┴───────────┴──────────────┘
```