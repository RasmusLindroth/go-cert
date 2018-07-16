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

//domainsHolder holds multiple DomainData for printing JSON
type domainsHolder struct {
	Domains []*domain.Data `json:"domains"`
}

//printTable prints a table with the domains
func printTable(conf *config, domains []*domain.Domain) {
	w := conf.WriterValue["out"]
	t := tb.InitTable(3, " ", []uint{tb.CenterHeader, tb.AlignRight | tb.CenterHeader, tb.CenterHeader, tb.AlignCenter | tb.CenterHeader})

	headerFormat := []tb.Attribute{}
	if conf.BoolValue["formatting"] == true {
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
		data := d.GetData(conf.LocationValue["location"])
		if conf.BoolValue["onlyExpiring"] && data.DaysLeft >= conf.IntValue["minDays"] {
			continue
		}
		dayColor := []tb.Attribute{}
		statusColor := []tb.Attribute{}

		if conf.BoolValue["colors"] {
			dayColor = []tb.Attribute{tb.FgGreen}
			statusColor = []tb.Attribute{tb.FgGreen}
		}
		if data.DaysLeft < conf.IntValue["minDays"] && conf.BoolValue["colors"] {
			dayColor[0] = tb.FgRed
		}
		if d.Error != nil && conf.BoolValue["colors"] {
			statusColor[0] = tb.FgRed
		}
		t.AddRow(
			[]*tb.Column{
				&tb.Column{Text: d.Name},
				&tb.Column{Text: strconv.Itoa(data.DaysLeft), Format: dayColor},
				&tb.Column{Text: data.EndTime.Format("2006-01-02 15:04")},
				&tb.Column{Text: data.Status, Format: statusColor},
			},
		)
	}
	t.Print(w)
}

func main() {
	loc := time.Local

	conf := &config{
		WriterValue: map[string]io.Writer{
			"out": os.Stdout,
		},

		IntValue: map[string]int{
			"minDays": 20,
		},

		StringValue: map[string]string{
			"outputType": "table",
		},

		StringValues: map[string][]string{
			"domains": []string{},
		},

		BoolValue: map[string]bool{
			"colors":       false,
			"formatting":   false,
			"onlyExpiring": false,
		},

		LocationValue: map[string]*time.Location{
			"location": loc,
		},
	}

	if checkConfigDefaults(conf) == false {
		log.Panic("You must set all default values for every item in your config struct")
	}

	app := cli.NewApp()
	app.Name = "go-cert"
	app.Usage = "check days left on SSL certificates"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Rasmus Lindroth",
			Email: "rasmus@lindroth.xyz",
		},
	}
	app.UsageText = "go-cert [OPTION]... DOMAIN [DOMAIN ...]"
	app.ArgsUsage = "domain [domain...]"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "days, d",
			Usage: "days `INT` left on certificate warning",
			Value: conf.IntValue["minDays"],
		},
		cli.StringFlag{
			Name:  "location, l",
			Usage: "`LOC` used for time zone, e.g. Europe/Stockholm. Defaults to local",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output `TYPE`: table, json, text (| seperator)",
			Value: conf.StringValue["outputType"],
		},
		cli.BoolFlag{
			Name:  "expiring, e",
			Usage: "only list certs where (days left < --days)",
		},
		cli.BoolFlag{
			Name:  "colors, c",
			Usage: "add colors in table output",
		},
		cli.BoolFlag{
			Name:  "formatting, f",
			Usage: "add bold in table header",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			for _, d := range c.Args() {
				conf.StringValues["domains"] = append(conf.StringValues["domains"], d)
			}
		}

		mapConfCli(conf, c)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	domains := []*domain.Domain{}

	for _, s := range conf.StringValues["domains"] {
		d, err := domain.InitDomain(s)

		if err != nil {
			//print error to stderr?
		}
		domains = append(domains, d)
	}

	if len(conf.StringValues["domains"]) == 0 {
		return
	}

	if conf.StringValue["outputType"] == "table" {
		printTable(conf, domains)
	} else if conf.StringValue["outputType"] == "json" {
		domainsData := domainsHolder{Domains: []*domain.Data{}}
		for _, d := range domains {
			data := d.GetData(conf.LocationValue["location"])
			if conf.BoolValue["onlyExpiring"] && data.DaysLeft >= conf.IntValue["minDays"] {
				continue
			}
			domainsData.Domains = append(domainsData.Domains, data)
		}
		jb, _ := json.Marshal(domainsData)
		fmt.Println(string(jb))
	} else if conf.StringValue["outputType"] == "text" {
		for _, d := range domains {
			data := d.GetData(conf.LocationValue["location"])
			if conf.BoolValue["onlyExpiring"] && data.DaysLeft >= conf.IntValue["minDays"] {
				continue
			}
			conf.WriterValue["out"].Write([]byte(
				fmt.Sprintf("%s|%d|%s|%s\n", data.Name, data.DaysLeft, data.EndTime, data.Status),
			))
		}
	}
}
