package main

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/republique/compiler"
)

func compile(log *logrus.Logger, args []string) error {
	c := compiler.NewCompiler(log)
	for _, v := range args {
		c.Compile(v)
	}
	return nil
}
