package main

import (
	"errors"
	"github.com/LeakIX/fmexplore/fmclient"
	"log"
	"os"
	"path"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("./fmexplore https://username:password@fminstance.com ./output-dir")
	}
	dirInfo, err := os.Stat(os.Args[2])
	if err == nil {
		log.Fatalln(errors.New("output directory already exists"))
	}
	err = os.Mkdir(os.Args[2], 0700)
	if err != nil {
		log.Fatalln(errors.New("couldn't create output directory"))
	}
	dirInfo, err = os.Stat(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	if !dirInfo.IsDir() {
		log.Fatalln(errors.New("couldn't create output directory"))
	}
	fmClient := fmclient.GetFmClient(os.Args[1])
	databases, err := fmClient.GetDatabases()
	if err != nil {
		log.Fatalln(err)
	}
	for _, database := range databases {
		log.Printf("Found database %s", database.Name)
		err = fmClient.AuthDatabase(database.Name)
		if err != nil {
			log.Fatalln(err)
		}
		layouts, err := fmClient.GetLayouts(database.Name)
		if err != nil {
			log.Fatalln(err)
		}
		for _, layout := range layouts {
			filename := path.Join(dirInfo.Name(), database.Name+"-"+layout.Name+".json")
			log.Printf("Dumping database %s, layout %s into %s", database.Name, layout.Name, filename)
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
