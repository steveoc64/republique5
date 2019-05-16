package republique

import "github.com/sirupsen/logrus"

type Compiler struct {
	log *logrus.Logger
}

func NewCompiler(log *logrus.Logger) *Compiler {
	return &Compiler{log: log}
}

func (c *Compiler) Compile(filename string) error {
	c.log.WithField("filename", filename).Print("Compiling")
	return nil
}
