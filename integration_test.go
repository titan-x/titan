// +build integration

package main

// var fooAddr = flag.String(...)
//
// func TestToo(t *testing.T) {
//     f, err := foo.Connect(*fooAddr)
//     // ...
// }

// go test takes build tags just like go build, so you can call go test -tags=integration. It also synthesizes a package main which calls flag.Parse,
// so any flags declared and visible will be processed and available to your tests.
//
// When you do this	Run this
// Save	go fmt (or goimports)
// Build	go vet, golint, and maybe go test
// Deploy	go test -tags=integration
//
// or is the gcm style t.Short() better?
