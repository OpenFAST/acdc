//go:build gen

package main

import "acdc/fio/schema"

func main() {
	schema.GenerateStructs("structs.go")
}
