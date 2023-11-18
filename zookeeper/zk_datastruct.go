package zookeeper

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/scalog/scalog/logger"
)

/* Trie data structure as the underlying implementation of ZK*/
type TrieNode struct {
	Value    string
	Children map[string]*TrieNode
	/*
	  Data and operation both are log structured aka each change
	  is logged increasing order
	*/
	data      []string
	operation []string
	version   int
}

type Trie struct {
	Root *TrieNode
}

// Upsert Inserts a path component into the trie. If exists - updates
func (t *Trie) Upsert(components []string, op string, data string) {
	current := t.Root
	for _, component := range components {
		if current.Children == nil {
			current.Children = make(map[string]*TrieNode)
		}
		if _, exists := current.Children[component]; !exists {
			current.Children[component] = &TrieNode{Value: component}
		}
		current = current.Children[component]
	}

	// Add the data to the log
	current.data = append(current.data, data)
	current.operation = append(current.operation, op)
	current.version += 1
}

// Get version, data and last operation from the Trie
func (t *Trie) Get(components []string, opts ...int) (ver int, op string, data string) {
	current := t.Root
	version := -1

	for _, component := range components {
		// Path does not exist: TODO find a way to return better error codes
		if _, exists := current.Children[component]; !exists {
			return -1, "false", "false"
		}
		current = current.Children[component]
	}

	if len(opts) >= 1 && opts[0] != -1 {
		version = opts[0]
	}

	if version < 0 || version >= current.version-1 {
		version = current.version - 1
	}

	return version, current.operation[version], current.data[version]
}

// Delete removes a path component from the trie.
func (t *Trie) Delete(components []string) bool {
	current := t.Root
	var parents []*TrieNode

	// Traverse the trie to find the node to delete and its parents
	for _, component := range components {
		if current.Children == nil {
			return false // Node not found
		}
		parents = append(parents, current)
		current = current.Children[component]
	}

	if current != nil {
		// Remove the node by deleting it from its parent's children
		for i := len(parents) - 1; i >= 0; i-- {
			parent := parents[i]
			delete(parent.Children, components[i])
			if len(parent.Children) > 0 {
				break
			}
		}
		current.version -= 1
		return true // Node successfully deleted
	}
	return false // Node not found
}

// PrintTrie prints the trie in a structured format.
func (t *Trie) PrintTrie(node *TrieNode, depth int) string {
	if node == nil {
		return ""
	}

	serialisedTrie := ""

	for i := 0; i < depth; i++ {
		serialisedTrie += "  "
		// log.Print("  ")
	}
	serialisedTrie += fmt.Sprintf("%s:: [", node.Value)
	// log.Printf("%s:: [", node.Value)

	sz := len(node.data) - 1
	if sz >= 0 {
		for i := 0; i < sz; i++ {
			serialisedTrie += fmt.Sprintf("%s::%s, ", node.operation[i], node.data[i])
			// log.Printf("%s::%s, ", node.operation[i], node.data[i])
		}
		serialisedTrie += fmt.Sprintf("%s::%s", node.operation[sz], node.data[sz])
		// log.Printf("%s::%s", node.operation[sz], node.data[sz])
	}
	serialisedTrie += fmt.Sprintf("]\n")
	// log.Printf("]\n")

	for _, child := range node.Children {
		serialisedTrie += t.PrintTrie(child, depth+1)
	}

	return serialisedTrie
}

/* Parsing of the main operation received from ZK client */
/* Returns array of split location, operations, and value */
func parse(input string) ([]string, string, string) {
	var operation string
	var location string
	var value string

	// Split the input string using "::" as the delimiter
	parts := strings.Split(input, "::")

	if len(parts) >= 2 {
		operation = parts[0]
		location = parts[1]
	}

	if len(parts) == 3 {
		value = parts[2]
	} else {
		value = "-1"
	}

	components := strings.Split(location, "/")
	components = components[1:]

	return components, operation, value
}

func (trie *Trie) Execute(data string) {
	pathComponents, operation, value := parse(data)

	switch operation {
	case "CREATE":
		trie.Upsert(pathComponents, operation, value)
		// trie.PrintTrie(trie.Root, 0)
	case "UPDATE":
		trie.Upsert(pathComponents, operation, value)
		// trie.PrintTrie(trie.Root, 0)
	case "DELETE_DATA":
		trie.Upsert(pathComponents, operation, " ")
		// trie.PrintTrie(trie.Root, 0)
	case "DELETE":
		ret := trie.Delete(pathComponents)
		if !ret {
			log.Printf("[ Zookeeper ][ Datastructure ][ Excute ][ DELETE ]Path `%s` does not exist.\n", strings.Join(pathComponents, "/"))
		}
	}

}

func (trie *Trie) ExecuteGet(command string) (int, string, error) {
	pathComponents, _, value := parse(command)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		err_ := fmt.Sprintf("[ Zookeeper ][ Datastructure ][ ExecuteGet ]Path `%v`. Error in version value parsing: %v", command, err)
		log.Printf(err_)
		return -1, "", err
	}

	ver, _, data := trie.Get(pathComponents, intValue)
	if ver == -1 {
		err_ := fmt.Sprintf("[ Zookeeper ][ Datastructure ][ ExecuteGet ]Path `%v`. Failed to get: %v", command, err)
		log.Printf(err_)
		return -1, "", err
	}

	return ver, data, nil
}

/* Only for testing - Use package main to test independently */
func main() {
	var operation string
	var value string
	var pathComponents []string

	bulk_commands := [...]string{
		"CREATE::/a/b/c/d::data1",
		"UPDATE::/a/b/c/d::data2",
		"CREATE::/a/b/c/e::data3",
		"CREATE::/a/b/c::data4",
		"GET::/a/b/c::0",
		"GET::/a/b/c/d::0",
		"GET::/a/b/c/d::1",
		"GET::/a/b/c/d::23",
		"GET::/a/b/c/e",
		"DELETE_DATA::/a/b/c/e",
		"DELETE::/a/b/c/e",
	}

	trie := &Trie{Root: &TrieNode{Value: "/"}}

	for _, command := range bulk_commands {
		pathComponents, operation, value = parse(command)
		if operation == "" {
			fmt.Println("Parse failed. Invalid input")
			os.Exit(0)
		}
		fmt.Println(command)

		switch operation {
		case "CREATE":
			trie.Upsert(pathComponents, operation, value)
			trie.PrintTrie(trie.Root, 0)
		case "UPDATE":
			trie.Upsert(pathComponents, operation, value)
			trie.PrintTrie(trie.Root, 0)
		case "GET":
			intValue, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Parse failed. Invalid input")
				os.Exit(0)
			}

			ver, op, data := trie.Get(pathComponents, intValue)
			if ver == -1 {
				fmt.Println("Failed to GET")
				os.Exit(0)
			}
			fmt.Printf("VER:%d :: %s :: %s\n", ver, op, data)
		case "DELETE_DATA":
			trie.Upsert(pathComponents, operation, " ")
			trie.PrintTrie(trie.Root, 0)
		case "DELETE":
			ret := trie.Delete(pathComponents)
			if !ret {
				fmt.Printf("Path `%s` does not exist.\n", strings.Join(pathComponents, "/"))
			}
			trie.PrintTrie(trie.Root, 0)
		}
		fmt.Println()
	}
}
