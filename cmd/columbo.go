package main

import (
	"github.com/battlemidget/columbo-go/internal/rules"
	"github.com/battlemidget/columbo-go/internal/tarextract"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func main() {
	app := &cli.App{
		Name:  "columbo",
		Usage: "he got clues for dayzzz",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "rules",
				Aliases: []string{"r"},
				Usage:   "Load rules spec from `YAML` file",
			},
			&cli.StringFlag{
				Name:    "outdir",
				Aliases: []string{"o"},
				Usage:   "Output `DIR` to store extracted files",
			},
		},
		Action: func(c *cli.Context) error {
			var rulesSpec rules.RulesSpec
			var srcFile string

			if c.NArg() > 0 {
				srcFile = c.Args().Get(0)
			}

			outDir := c.String("outdir")
			err := tarextract.Extract(srcFile, outDir)
			if err != nil {
				log.Fatal(err)
			}
			// glob extracted file looking for any other compressed files
			// basically a max depth 1 search for now
			matches, _ := filepath.Glob(filepath.Join(outDir, "*tar*"))
			for _, match := range matches {
				log.Println("Found ", match)
				err = tarextract.Extract(match, outDir)
				if err != nil {
					log.Println("Could not extract, ", err)
					continue
				}
			}
			parseRules := c.String("rules")
			output := rulesSpec.Parse(parseRules)
			for _, r := range output.Rules {
				if r.LineMatch != "" {
					log.Println("Processing: ", r.Id)
					r.ProcessLineMatch(outDir)
				}

				if r.StartMarker != "" && r.EndMarker != "" {
					log.Println("Processing: ", r.Id)
					r.ProcessStartEndMarker(outDir)
				}
			}
			rules.SaveResults()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
