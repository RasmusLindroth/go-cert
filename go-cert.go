package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/RasmusLindroth/go-cert/domain"
	tb "github.com/RasmusLindroth/go-cert/table"
	"github.com/urfave/cli"
)

//Config holds all config params
type Config struct {
	//MinDays, when to warn for cert expiration
	Out        io.Writer
	MinDays    int
	Domains    []string
	PrintTable bool
	PrintJSON  bool
	Colors     bool
	Formatting bool
	Location   *time.Location
}

//Domains holds multiple DomainData for printing JSON
type Domains struct {
	Domains []*domain.DomainData `json:"domains"`
}

//PrintTable prints a table with the domains
func PrintTable(conf *Config, w io.Writer, domains []*domain.Domain) {
	t := tb.InitTable(3, " ", []uint{tb.CenterHeader, tb.AlignRight | tb.CenterHeader, tb.CenterHeader, tb.AlignCenter | tb.CenterHeader})

	headerFormat := []tb.Attribute{}
	if conf.Formatting == true {
		headerFormat = append(headerFormat, tb.Bold)
	}

	t.AddRow(
		[]*tb.Column{
			&tb.Column{Text: "Domain", Format: headerFormat},
			&tb.Column{Text: "Days left", Format: headerFormat},
			&tb.Column{Text: "End date", Format: headerFormat},
			&tb.Column{Text: "Status", Format: headerFormat},
		},
	)
	for _, d := range domains {
		domainData := d.GetData(conf.Location)
		dayColor := []tb.Attribute{}
		statusColor := []tb.Attribute{}

		if conf.Colors {
			dayColor = []tb.Attribute{tb.FgGreen}
			statusColor = []tb.Attribute{tb.FgGreen}
		}
		if domainData.DaysLeft < conf.MinDays && conf.Colors {
			dayColor[0] = tb.FgRed
		}
		if d.Error != nil && conf.Colors {
			statusColor[0] = tb.FgRed
		}
		t.AddRow(
			[]*tb.Column{
				&tb.Column{Text: d.Name},
				&tb.Column{Text: strconv.Itoa(domainData.DaysLeft), Format: dayColor},
				&tb.Column{Text: domainData.EndTime.Format("2006-01-02 15:04")},
				&tb.Column{Text: domainData.Status, Format: statusColor},
			},
		)
	}
	t.Print(w)
}

func main() {
	loc := time.Local
	conf := &Config{
		Out:        os.Stdout,
		MinDays:    20,
		Domains:    []string{},
		PrintTable: false,
		PrintJSON:  false,
		Colors:     false,
		Formatting: false,
		Location:   loc,
	}

	app := cli.NewApp()
	app.Name = "go-cert"
	app.Usage = "check days left on SSL certificates"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Rasmus Lindroth",
			Email: "rasmus@lindroth.xyz",
		},
	}
	app.ArgsUsage = "domains"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "days, d",
			Value: 20,
			Usage: "days left on certificate warning",
		},
		cli.StringFlag{
			Name:  "location, l",
			Usage: "used for time zone, e.g. Europe/Stockholm. Defaults to local",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "table (default), json",
		},
		cli.BoolFlag{
			Name:  "colors, c",
			Usage: "add colors in table output",
		},
		cli.BoolFlag{
			Name:  "formatting, f",
			Usage: "uses bold in table header",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			for _, d := range c.Args() {
				conf.Domains = append(conf.Domains, d)
			}
		}

		conf.MinDays = c.Int("days")
		if c.String("location") != "" {
			loc, err := time.LoadLocation(c.String("location"))
			if err == nil {
				conf.Location = loc
			}
		}

		if c.String("output") != "table" && c.String("output") != "json" {
			//Print error because not valid?
			conf.PrintTable = true
			conf.PrintJSON = false
		} else if c.String("output") == "table" {
			conf.PrintTable = true
			conf.PrintJSON = false
		} else if c.String("output") == "json" {
			conf.PrintJSON = true
			conf.PrintTable = false
		}

		if c.Bool("colors") {
			conf.Colors = c.Bool("colors")
		}
		if c.Bool("formatting") {
			conf.Formatting = c.Bool("formatting")
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	domains := []*domain.Domain{}

	for _, s := range conf.Domains {
		d, err := domain.InitDomain(s)

		if err != nil {
			//print error to stderr?
		}
		domains = append(domains, d)
	}

	if conf.PrintTable {
		PrintTable(conf, conf.Out, domains)
	} else if conf.PrintJSON {
		domainsData := Domains{Domains: []*domain.DomainData{}}
		for _, d := range domains {
			domainsData.Domains = append(domainsData.Domains, d.GetData(conf.Location))
		}
		jb, _ := json.Marshal(domainsData)
		fmt.Println(string(jb))
	}
}
