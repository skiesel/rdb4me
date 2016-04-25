package main

import (
	"bufio"
	"fmt"
	"os"
	"io/ioutil"
	"strings"

	"github.com/skiesel/expsys/rdb"
)

var (
	root = ""
	filters = map[string]string{}
	ds *rdb.Dataset
)

func main() {
	var err error
	root, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("RDB-4-Me!")
	fmt.Printf("current root set to: %s\n", root)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if handleCommand(reader, text) {
			break
		}
	}
}

func handleCommand(stdin *bufio.Reader, text string) bool {
	tokens := strings.Fields(text)
	if len(tokens) <= 0 {
		return false
	}

	command := tokens[0]

	switch command {
	case "quit":
		return true
	
	case "set-root":
		root = strings.Join(tokens[1:], " ")
		fmt.Printf("current root set to: %s\n", root)

	case "print-root":
		fmt.Printf("current root set to: %s\n", root)

	case "print-filter":
		for key, val := range filters {
			fmt.Printf("%s = %s\n", key, val)
		}

	case "add-filter":
		filter := strings.Split(strings.Join(tokens[1:], " "), "=")

		if len(filter) < 2 {
			fmt.Printf("ignoring malformed filter add: %s\n", strings.Join(filter[1:], "="))
		} else {
			key := filter[0]
			val := filter[1]
			filters[key] = val
		}

	case "read-filter":
		readFilter(tokens)

	case "read-dataset":
		ds = rdb.GetDatasetWithPathKeys(root, filters, "ActiveDataset")

	case "dataset-size":
		if ds != nil {
			fmt.Println(ds.GetSize())
		} else {
			fmt.Println("no active dataset")
		}

	case "delete-dataset-files":
		if ds != nil {
			fmt.Printf("This *WILL* delete actively selected %d files from the FS. Are you sure? (yes/no)\n", ds.GetSize())
			yesOrNo, err := stdin.ReadString('\n')
			yesOrNo = strings.TrimSpace(yesOrNo)
			if yesOrNo != "yes" || err != nil {
				fmt.Println("not deleted")
			} else {
				fmt.Println("deleting")
				deletedDatasetFiles()
			}
		} else {
			fmt.Println("no active dataset")
		}

	// case "rename-path-key":
	// 	if ds != nil {
	// 		if len(tokens) < 4 {
	// 			fmt.Println("ignoring malformed rename request")
	// 		} else {
	// 			renamePathKey(tokens[1], tokens[2], tokens[3])
	// 		}
	// 	} else {
	// 		fmt.Println("no active dataset")
	// 	}

	case "dataset-files":
		if ds != nil {
			paths := ds.GetDatasetPathes()
			for i, path := range paths {
				if i >= 25 {
					fmt.Printf("... %d more\n", ds.GetSize() - i)
					break
				}
				fmt.Println(path)
			}
		} else {
			fmt.Println("no active dataset")
		}		

	default:
		fmt.Printf("unrecognized command: %s\n", command)
	}

	return false
}

func readFilter(tokens []string) {
	if len(tokens) < 2 {
		return
	}

	file, err := os.Open(tokens[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	filters = map[string]string{}

	reader := bufio.NewReader(file)
	
	for {
		text, err := reader.ReadString('\n')
		tokens := strings.Split(text, "=")
		if len(tokens) < 2 {
			fmt.Printf("skipping malformed filter: %s\n", text)
			if err != nil {
				break
			} else {
				continue
			}
		}
		key := strings.TrimSpace(tokens[0])
		val := strings.TrimSpace(tokens[1])
		filters[key] = val
		if err != nil {
			break
		}
	}
}

func deletedDatasetFiles() {
	paths := ds.GetDatasetPathes()
	if len(paths) <= 0 {
		return
	}

	ds = nil
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getKeyInDirectory(directory string) string {
	fInfo, error := ioutil.ReadDir(directory)

	if error != nil { // there was an error
		panic(error)
	} else {
		for i := range fInfo {
			if strings.Contains(fInfo[i].Name(), "KEY=") {
				return strings.SplitAfter(fInfo[i].Name(), "KEY=")[1]
			}
		}
	}
	fmt.Println("No key file found in ", directory)
	panic("No key file found")
}