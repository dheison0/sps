package forwards

import (
	"log"
	"net/url"
	"regexp"
	"sps/types"
	"sps/util"
)

var FilterRegex = []*regexp.Regexp{}
var Filter = map[string]*string{}
var config types.FilterConfig

func SetConfigAndParse(c types.FilterConfig) {
	config = c
	if c.File == "" {
		return
	}
	line, err := ParseFilterFile()
	if err != nil {
		log.Fatalf("Error on line %d: %s\n", line, err.Error())
	}
}

func ParseFilterFile() (int, error) {
	log.Println("Processing filter file...")
	data, err := util.ReadFile(config.File)
	if err != nil {
		return 0, err
	}
	lines := util.ReadLinesFromBytes(data)
	for i, l := range lines {
		line := string(l)
		Filter[line] = nil
		if config.EnableRegex == false {
			continue
		}
		if config.LessMemory {
			Filter[line] = nil
		} else {
			re, err := regexp.Compile(line)
			if err != nil {
				return i + 1, err
			}
			FilterRegex = append(FilterRegex, re)
		}
	}
	log.Println("Filter file proccess finished!")
	return len(lines), nil
}

func MatchFilter(matchURL string) bool {
	log.Println("Testing filters...")
	result := false
	urlInfo, _ := url.Parse(matchURL)
	if _, ok := Filter[urlInfo.Host]; ok {
		result = true
	} else if config.EnableRegex && config.LessMemory {
		for pattern, _ := range Filter {
			match, err := regexp.MatchString(pattern, matchURL)
			if err != nil {
				log.Println("Pattern '%s' compilation error: %s", pattern, err.Error())
				continue
			}
			if match {
				result = true
				break
			}
		}
	} else if config.EnableRegex {
		for _, regex := range FilterRegex {
			match := regex.MatchString(matchURL)
			if match {
				result = true
				break
			}
		}
	}
	if result {
		log.Printf("%s blocked!\n", matchURL)
	} else {
		log.Println("Done, alright!")
	}
	return result
}
