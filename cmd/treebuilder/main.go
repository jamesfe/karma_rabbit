package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type CommentTree struct {
	Item       Comment `json:"comment"`
	parent_obj *CommentTree
	Children   []*CommentTree `json:"children"`
}

type Comment struct {
	Score     int    `json:"score"`
	Comment   string `json:"body"`
	Subreddit string `json:"subreddit"`
	Parent    string `json:"parent_id"`
	Id        string `json:"id"`
	CreatedTS int64  `json:"created_utc"`
}

func file_exists(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	return true
}

func InsertToTree(tree *CommentTree, comment Comment) bool {
	comment_parent := strings.Split(comment.Parent, "_")[1]
	//	fmt.Printf("Parent: %s Child: %s Child Unparsed %s \n", tree.item.Id, comment_parent, comment.Parent)
	if tree.Item.Id == comment_parent {
		tree.Children = append(tree.Children, &CommentTree{Item: comment, parent_obj: tree, Children: []*CommentTree{}})
		fmt.Printf("appending something!\n")
		fmt.Printf("Children Count: %d\n", len(tree.Children))
		return true
	} else {
		for _, subitem := range tree.Children {
			val := InsertToTree(subitem, comment)
			if val == true {
				return true
			}
		}
	}
	return false
}

func main() {
	// We have a list of comments and a list of parent comments
	var t3Comments []CommentTree
	var commentBag []Comment

	if len(os.Args) < 3 {
		fmt.Printf("Use this program as such: treebuilderinput_file output_folder")
		os.Exit(-1)
	}
	target_file := os.Args[1]
	output_folder := os.Args[2]
	if !file_exists(target_file) {
		fmt.Printf("Target file does not exist.")
		os.Exit(-1)
	}
	if !file_exists(output_folder) {
		fmt.Printf("Output folder does not exist.")
		os.Exit(-1)
	}
	// Read the comment file
	file, err := os.Open(target_file)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var dat Comment
	for scanner.Scan() {
		input := scanner.Bytes()
		err = json.Unmarshal(input, &dat)
		commentBag = append(commentBag, dat)
	}
	// By sorting the comments we can be sure it will only take one pass to build the tree
	sort.Slice(commentBag, func(i, j int) bool {
		return commentBag[i].CreatedTS < commentBag[j].CreatedTS
	})

	// Pass through the list, adding t1 items and filtering in other items
	for _, comment := range commentBag {
		if strings.HasPrefix(comment.Parent, "t3") {
			// Dump the comment into the list
			newTreeItem := CommentTree{Item: comment, parent_obj: nil, Children: []*CommentTree{}}
			t3Comments = append(t3Comments, newTreeItem)
			// fmt.Printf("found a t3 with date: %d \n", comment.CreatedTS)
		} else {
			for index, _ := range t3Comments {
				tree := &t3Comments[index]
				if InsertToTree(tree, comment) == true {
					fmt.Printf("Length: %d\n", len(tree.Children))
				}
			}
		}
	}
	fmt.Printf("checking for comments in %d items...\n\n", len(t3Comments))
	for _, c := range t3Comments {
		if len(c.Children) > 0 {
			fmt.Printf("Children: %d\n", len(c.Children))
		}
	}
	bt, err := json.Marshal(t3Comments)
	fmt.Printf("%s", bt)

}
