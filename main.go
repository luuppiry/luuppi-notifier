package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/luuppiry/luuppi-rss-service/fetchers"
	"github.com/luuppiry/luuppi-rss-service/formatters"
	"github.com/luuppiry/luuppi-rss-service/output"
)

var configPath = flag.String("configPath", "/config.json", "path to config.json")
var serverPort = flag.Uint("port", 42069, "Port where data is served")

type Config struct {
	News []News
}

type News struct {
	Source    Component
	Formatter Component
	Output    Component
}
type Component struct {
	ComponentType string
	Conf          map[string]string
}
type Formatter interface {
	Format(data any) ([]byte, error)
}

func ChooseFormatter(output string, input string, conf map[string]string) (Formatter, error) {
	spec := fmt.Sprintf("%s-%s", input, output)
	switch spec {
	case "news-rss":
		return formatters.NewRssNewsFormatter(conf), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown formatter: %s for %s data", output, input))
	}
}

type Fetcher interface {
	Fetch() (any, error)
}

func ChooseFetcher(sourceType string, fetcherType string, conf map[string]string) (Fetcher, error) {
	spec := fmt.Sprintf("%s-%s", sourceType, fetcherType)
	switch spec {
	case "strapiv4-news":
		return fetchers.NewStrapiv4NewsFetcher(conf), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown fetcher: %s for: %s", sourceType, fetcherType))
	}
}

type Outputter interface {
	Initialize() error
	Update([]byte) error
}

func ChooseOutputter(outputType string, conf map[string]string) (Outputter, error) {
	switch outputType {
	case "http":
		return output.NewHttpOutput(conf), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown outputter: %s", outputType))
	}
}

func ParseConfig(data []byte) Config {
	c := Config{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		log.Fatal("Failed to parse config", err)
	}
	return c
}

func ReadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to read config", err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Failed to read config", err)
	}
	c := ParseConfig(data)
	return c
}

func pipeline(s Fetcher, f Formatter, o Outputter) {
	ticker := time.NewTicker(1 * time.Minute)

	for _ = range ticker.C {
		data, err := s.Fetch()
		if err != nil {
			log.Printf("Fetching data failed: %s", err)
			continue
		}
		out, err := f.Format(data)
		if err != nil {
			log.Printf("Formatting data failed: %s", err)
			continue
		}
		err = o.Update(out)
		if err != nil {
			log.Printf("Outputting data failed: %s", err)
			continue
		}
	}
}

func initialize(conf Config) error {
	for _, n := range conf.News {
		source, err := ChooseFetcher(n.Source.ComponentType, "news", n.Source.Conf)
		if err != nil {
			return err
		}
		formatter, err := ChooseFormatter(n.Formatter.ComponentType, "news", n.Formatter.Conf)
		if err != nil {
			return err
		}
		outputter, err := ChooseOutputter(n.Output.ComponentType, n.Output.Conf)
		if err != nil {
			return err
		}
		outputter.Initialize()
		go pipeline(source, formatter, outputter)

	}
	return nil
}

func main() {
	flag.Parse()
	conf := ReadConfig(*configPath)
	initialize(conf)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil))

}
