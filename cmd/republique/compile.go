package main

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/republique"
)

func compile(log *logrus.Logger, args []string) error {
	c := republique.NewCompiler(log)
	for _, v := range args {
		c.Compile(v)
	}
	return nil
}
