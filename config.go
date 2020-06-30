package main

import (
	"io"
	"time"

	"github.com/urfave/cli/v2"
)

type config struct {
	WriterValue   map[string]io.Writer
	IntValue      map[string]int
	StringValue   map[string]string
	StringValues  map[string][]string
	BoolValue     map[string]bool
	LocationValue map[string]*time.Location
}

func checkConfigDefaults(conf *config) bool {
	writerKeys := []string{"out"}
	intKeys := []string{"minDays"}
	stringKeys := []string{"outputType"}
	stringMultiKeys := []string{"domains"}
	boolKeys := []string{"colors", "formatting", "onlyExpiring"}
	locationKeys := []string{"location"}

	for _, k := range writerKeys {
		if _, found := conf.WriterValue[k]; found == false {
			return false
		}
	}

	for _, k := range intKeys {
		if _, found := conf.IntValue[k]; found == false {
			return false
		}
	}

	for _, k := range stringKeys {
		if _, found := conf.StringValue[k]; found == false {
			return false
		}
	}

	for _, k := range stringMultiKeys {
		if _, found := conf.StringValues[k]; found == false {
			return false
		}
	}

	for _, k := range boolKeys {
		if _, found := conf.BoolValue[k]; found == false {
			return false
		}
	}

	for _, k := range locationKeys {
		if _, found := conf.LocationValue[k]; found == false {
			return false
		}
	}

	return true
}

func mapConfCli(conf *config, c *cli.Context) *config {
	configMappig := []unknownValue{
		&boolValue{
			FlagConf: &flagConfMap{
				FlagName:   "colors",
				ConfigName: "colors",
			},
		},
		&boolValue{
			FlagConf: &flagConfMap{
				FlagName:   "formatting",
				ConfigName: "formatting",
			},
		},
		&boolValue{
			FlagConf: &flagConfMap{
				FlagName:   "expiring",
				ConfigName: "onlyExpiring",
			},
		},
		&intValue{
			FlagConf: &flagConfMap{
				FlagName:   "days",
				ConfigName: "minDays",
			},
		},
		&stringValue{
			FlagConf: &flagConfMap{
				FlagName:   "output",
				ConfigName: "outputType",
			},
			Valid: []string{
				"table",
				"json",
				"text",
			},
		},
		&locationValue{
			FlagConf: &flagConfMap{
				FlagName:   "location",
				ConfigName: "location",
			},
		},
	}

	for _, item := range configMappig {
		if item.isSet(c) && item.isValid(c) {
			item.setValue(conf, c)
		}
	}

	return conf
}

type unknownValue interface {
	isSet(c *cli.Context) bool
	isValid(c *cli.Context) bool
	setValue(conf *config, c *cli.Context)
}

type boolValue struct {
	FlagConf *flagConfMap
}

func (v *boolValue) isSet(c *cli.Context) bool {
	return c.IsSet(v.FlagConf.FlagName)
}

func (v *boolValue) isValid(c *cli.Context) bool {
	return true
}

func (v *boolValue) setValue(conf *config, c *cli.Context) {
	conf.BoolValue[v.FlagConf.ConfigName] = c.Bool(v.FlagConf.FlagName)
}

type intValue struct {
	FlagConf *flagConfMap
}

func (v *intValue) isSet(c *cli.Context) bool {
	return c.IsSet(v.FlagConf.FlagName)
}

func (v *intValue) isValid(c *cli.Context) bool {
	return true
}

func (v *intValue) setValue(conf *config, c *cli.Context) {
	conf.IntValue[v.FlagConf.ConfigName] = c.Int(v.FlagConf.FlagName)
}

type locationValue struct {
	FlagConf *flagConfMap
}

func (v *locationValue) isSet(c *cli.Context) bool {
	return c.IsSet(v.FlagConf.FlagName)
}

func (v *locationValue) isValid(c *cli.Context) bool {
	_, err := time.LoadLocation(c.String(v.FlagConf.FlagName))
	return err == nil
}

func (v *locationValue) setValue(conf *config, c *cli.Context) {
	l, _ := time.LoadLocation(c.String(v.FlagConf.FlagName))
	conf.LocationValue[v.FlagConf.ConfigName] = l
}

type stringValue struct {
	FlagConf *flagConfMap
	Valid    []string
}

func (v *stringValue) isSet(c *cli.Context) bool {
	return c.IsSet(v.FlagConf.FlagName)
}

func (v *stringValue) isValid(c *cli.Context) bool {
	if len(v.Valid) == 0 {
		return true
	}

	value := c.String(v.FlagConf.FlagName)
	valid := false

	for _, x := range v.Valid {
		if x == value {
			valid = true
			break
		}
	}

	return valid
}

func (v *stringValue) setValue(conf *config, c *cli.Context) {
	conf.StringValue[v.FlagConf.ConfigName] = c.String(v.FlagConf.FlagName)
}

type flagConfMap struct {
	FlagName   string
	ConfigName string
}
