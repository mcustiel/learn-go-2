// duplicate-folders-finder project duplicate-folders-finder.go
package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var root Node = ScanDirToTree("C:\\test")
	excluded := make(map[string]bool, 0)
	orderedFiles := make(map[int]string, 0)
	equals := make(map[string][]string)
	current := 0
	bfs(root, func(value string) {
		fmt.Println("Saving directory ", value, " with index ", current)
		orderedFiles[current] = value
		current++
	})

	for index := 0; index < len(orderedFiles); index++ {
		dirName := orderedFiles[index]
		fmt.Println("Looking for duplicates of ", dirName)
		if !isExcluded(dirName, excluded) {
			for _, otherDirName := range orderedFiles {
				if dirName != otherDirName && !isExcluded(otherDirName, excluded) {
					fmt.Println("-- Comparing with", otherDirName)
					if compareDirContents(GetDirectoryContents(dirName), GetDirectoryContents(otherDirName)) {
						if _, has := equals[dirName]; !has {
							fmt.Println("-- Creating the key and list")
							equals[dirName] = make([]string, 0)
						}
						fmt.Println("-- Adding duplicated to the list")
						equals[dirName] = addToList(equals[dirName], otherDirName)
						fmt.Println("-- Excluding duplicated from future searches")
						excluded[otherDirName] = true
					}
				}
			}
		} else {
			fmt.Println("It was already excluded")
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

func compareDirContents(dir1 []os.FileInfo, dir2 []os.FileInfo) bool {
	if len(dir1) == 0 {
		return false
	}
	matches := 0
	for _, fileInfo1 := range dir1 {
		for _, fileInfo2 := range dir2 {
			if fileInfo1.Name() == fileInfo2.Name() && fileInfo1.IsDir() == fileInfo2.IsDir() {
				matches++
			}
		}
	}
	fmt.Println(float32(matches) / float32(len(dir1)))
	return matches == len(dir1)
}

func addToList(stringList []string, element string) []string {
	fmt.Println("Extending slice")
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

func bfs(start Node, callback func(string)) {
	visited := make(map[string]bool, 0)
	walkTree(start, func(value string) {
		fmt.Println("Setting ", value, " to not visited")
		visited[value] = false
	})

	queue := list.New()
	queue.PushBack(start)
	for queue.Len() > 0 {
		current := queue.Front().Value.(Node)
		queue.Remove(queue.Front())
		for e := current.children.Front(); e != nil; e = e.Next() {
			currentNode := e.Value.(Node)
			if passed := visited[currentNode.value]; !passed {
				visited[currentNode.value] = true
				queue.PushBack(currentNode)
				callback(currentNode.value)
			}
		}
	}
}

func walkTree(node Node, callback func(string)) {
	callback(node.value)
	for e := node.children.Front(); e != nil; e = e.Next() {
		walkTree(e.Value.(Node), callback)
	}
}

func ScanDirToTree(dirName string) Node {
	graph := NewNode(dirName)

	for _, file := range GetDirectoryContents(dirName) {
		if file.IsDir() {
			graph.children.PushBack(ScanDirToTree(dirName + string(os.PathSeparator) + file.Name()))
		}
	}
	return graph
}

func GetDirectoryContents(dirName string) []os.FileInfo {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	return files
}
