package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func list(log *logrus.Logger, options []string) error {

	full := false
	short := true
	filter := ""
	for _, v := range options {
		switch v {
		case "--players":
			short = false
		case "--full":
			full = true
			short = false
		case "--help":
			println("List options:")
			println(" --players   Print player and login details")
			println(" --full      Print full game details")
			return nil
		default:
			filter = strings.ToLower(v)
		}
	}
	root := filepath.Join(os.Getenv("HOME"), ".republique")
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if strings.HasSuffix(path, ".db") {
				if filter != "" && !strings.Contains(strings.ToLower(path), filter) {
					return nil
				}
				if !short {
					println("===========================================================================")
				}
				info(log, f.Name(), full, short)
			}
		}
		return nil
	})
	if err != nil {
		println("Error:", err.Error())
	}
	return nil
}
