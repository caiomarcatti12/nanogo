package main

import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/nanogo"
)

func main() {
	nanogo.Bootstrap()

	environment, _ := di.Get[env.IEnv]()

	test := environment.GetEnv("test")

	println(test)
}
