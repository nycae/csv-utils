# CSV utils

This is a collection of utilities I use in order to operate CSV files.

## Encoder

The first helper is  called csv.Encoder. This is an interface closer to the json/yaml/xml encoder from the standard
library of Go.

This encoder only works with slices of structs with primitive type fields. Another restriction is that the value (and
not the pointer) of the slice should be passed to the "Encode" method. The type of the slice, on the other hand can be
both "whole type" and "pointer to type".

It also allows to add a "csv" to a field in order to override the default name.
A canonical example should be:

```go
func example() {
    c := []struct{
        Name string `csv:"name"`
        ID   int    `csv:"id"`
    }{{"A", 1}, {"B", 2}, {"C", 3}}
    csv.NewEncoder(os.Stdout).Encode(c)
}
```