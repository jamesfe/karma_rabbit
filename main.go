package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Comment struct {
	score     int    `json:"score"`
	comment   string `json:"body"`
	subreddit string `json:"subreddit"`
	parent    string `json:"parent_id"`
	id        string `json:"id"`
}

var dat map[string]interface{}
var target_reddits = [...]string{"funny", "AskReddit", "todayilearned", "worldnews", "Science", "pics", "gaming", "IAmA", "videos", "movies", "aww", "Music", "blog", "gifs", "news", "explainlikeimfive", "askscience", "EarthPorn", "books", "television", "mildlyinteresting", "Showerthoughts", "LifeProTips", "space", "DIY", "Jokes", "gadgets", "nottheonion", "sports", "food", "tifu"}

func inArray(target string, array []string) bool {
	for _, item := range array {
		if strings.ToLower(item) == target {
			return true
		}
	}
	return false
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	target_file := "/Users/j.ferrara/PersCode/reddit_classification/data/sample_100k.json"
	output_folder := "/Users/j.ferrara/PersCode/reddit_classification/data/subreddits/"
	outFileMap := make(map[string]*bufio.Writer)
	file, err := os.Open(target_file)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range target_reddits {
		temp_item, err := os.Create(output_folder + strings.ToLower(item) + ".json")
		check(err)
		defer temp_item.Close()
		outFileMap[strings.ToLower(item)] = bufio.NewWriter(temp_item)
	}
	// fmt.Printf("%#v", outFileMap)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &dat)
		sub := strings.ToLower(dat["subreddit"].(string))
		if inArray(sub, target_reddits[:]) {
			// fmt.Printf("%s", sub)
			outFileMap[sub].Write(scanner.Bytes())
			outFileMap[sub].WriteString("\n")
		}

	}
	for _, v := range outFileMap {
		v.Flush()
	}
}
