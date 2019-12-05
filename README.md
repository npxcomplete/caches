A slow growing collection of cache algorithms in golang 
and compatible with the genny generics generation tool.
(see examples/generate.go)

Note the use of wrapper types, and the lack of a publicly exported data type.
This is intentional by the maintainers. It is idiomatic in go for consumers to define
the interface of the types they consume, so if for some reason the contract 
of this library changes, that should not impact the expected behavior of 3rd party code.
Though not all of that contract can be communicated in the go type system we wish to get
as much bang for our buck as we can. Additionally the presence of a general cache interface
marginally increases the complexity of using the library. Similarly generation of the 
full implementation has some annoying edge cases where referencing non-generic code such 
as error values results in output that will not compile. Hence, generator targets are 
kept as simple as possible.   