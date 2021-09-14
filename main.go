package main

import (
	"github.com/LeakIX/fmexplore/fmclient"
	"log"
	"os"
)

func main() {
	fmClient := fmclient.GetFmClient(os.Args[1])
	databases, err := fmClient.GetDatabases()
	if err != nil {
		log.Fatalln(err)
	}
	for _, database := range databases {
		log.Println(database.Name)
		err = fmClient.AuthDatabase(database.Name)
		if err != nil {
			log.Fatalln(err)
		}
		layouts, err := fmClient.GetLayouts(database.Name)
		if err != nil {
			log.Fatalln(err)
		}
		for _, layout := range layouts {
			filename := database.Name + "-" + layout.Name + ".json"
			log.Printf("Dumping %s %s into %s", database.Name, layout.Name, filename)
			outputFile, err := os.Create(filename)
			if err != nil {
				log.Fatalln(err)
			}
			err = fmClient.Dump(database.Name, layout.Name, outputFile)
			if err != nil {
				log.Println(err)
			}
			err = outputFile.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

