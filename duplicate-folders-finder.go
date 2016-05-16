// duplicate-folders-finder project duplicate-folders-finder.go
package main

import (
	// "container/list"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mcustiel/graph"
)

func main() {
	var graphInstance *graph.Graph = ScanDirToTree("C:\\test")
	var equalityIndex float32 = 0.8

	excluded := make(map[string]bool, 0)
	orderedFiles := make(map[int]string, 0)
	equals := make(map[string][]string)

	current := 0
	graphInstance.Bfs(func(node *graph.Node) {
		value := node.Value().(string)
		orderedFiles[current] = value
		current++
	})

	for index := 0; index < len(orderedFiles); index++ {
		dirName := orderedFiles[index]
		if !isExcluded(dirName, excluded) {
			for _, otherDirName := range orderedFiles {
				if dirName != otherDirName && !isExcluded(otherDirName, excluded) {
					if calculateEqualityIndex(GetDirectoryContents(dirName), GetDirectoryContents(otherDirName)) >= equalityIndex {
						if _, has := equals[dirName]; !has {
							equals[dirName] = make([]string, 0)
						}
						equals[dirName] = addToList(equals[dirName], otherDirName)
						excluded[otherDirName] = true
					}
				}
			}
		}
	}

	for dirName, duplicates := range equals {
		fmt.Println("Duplicates for ", dirName)
		for _, duplicate := range duplicates {
			fmt.Println("    ", duplicate)
		}
	}
}

func isExcluded(dirName string, excluded map[string]bool) bool {
	for excludedDir, _ := range excluded {
		if strings.HasPrefix(dirName, excludedDir) {
			return true
		}
	}
	return false
}

func calculateEqualityIndex(dir1 []os.FileInfo, dir2 []os.FileInfo) float32 {
	if len(dir1) == 0 || len(dir2) == 0 {
		return 0
	}
	matches := 0
	for _, fileInfo1 := range dir1 {
		for _, fileInfo2 := range dir2 {
			if fileInfo1.Name() == fileInfo2.Name() && fileInfo1.IsDir() == fileInfo2.IsDir() {
				matches++
			}
		}
	}
	return float32(matches) / float32(len(dir1))
}

func addToList(stringList []string, element string) []string {
	n := len(stringList)
	if n == cap(stringList) {
		newList := make([]string, len(stringList), cap(stringList)+1)
		copy(newList, stringList)
		stringList = newList
	}
	stringList = stringList[0 : n+1]
	stringList[n] = element
	return stringList
}

func ScanDirToTree(dirName string) *graph.Graph {
	var function func(string) *graph.Node
	function = func(curDir string) *graph.Node {
		node := graph.NewNode(curDir)
		for _, file := range GetDirectoryContents(curDir) {
			if file.IsDir() {
				node.AddAdyacent(function(curDir + string(os.PathSeparator) + file.Name()))
			}
		}
		return node
	}
	return graph.New(function(dirName))
}

func GetDirectoryContents(dirName string) []os.FileInfo {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	return files
}
