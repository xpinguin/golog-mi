package util

import (
	"go/build"
	"log"

	"golang.org/x/tools/go/loader"
)

func LoadProgram(mainPkg string) *loader.Program {
	ldrCfg := loader.Config{
		Build:       &build.Default,
		AllowErrors: true,
	}
	ldrCfg.Import(mainPkg)

	prg, err := ldrCfg.Load()
	if err != nil {
		log.Fatal("Unable to load pkg: ", err)
		return nil
	}
	return prg
}
