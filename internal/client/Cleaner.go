package client

import (
	"log"
	"os"
)

type Cleaner struct {
}

func NewCleaner() *Cleaner {
	return &Cleaner{}
}

func (c *Cleaner) Clean(filenames []string) {
	for _, filename := range filenames {
		log.Printf("Removing file: %s\n", filename)
		err := os.Remove(filename)
		if err != nil {
			log.Printf("Cannot remove file: %s", err)
		}
	}
}
