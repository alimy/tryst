module github.com/alimy/tryst

go 1.24.0

require (
	github.com/RoaringBitmap/roaring v1.9.4
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/stretchr/testify v1.11.1
	golang.org/x/crypto v0.48.0
	golang.org/x/sync v0.19.0
	golang.org/x/sys v0.41.0
)

require (
	github.com/bits-and-blooms/bitset v1.12.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v1.20.0 // invalid version
