module hello

go 1.16

replace jokes => ../jokes

require (
	cache v0.0.0-00010101000000-000000000000
	jokes v0.0.0-00010101000000-000000000000
)

replace cache => ../cache
