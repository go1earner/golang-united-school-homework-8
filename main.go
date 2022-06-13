package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

const filePerm = 0600

type Item struct {
	Id, Email string
	Age       int
}

func add(item, filename string, writer io.Writer) {
	id, err := parseIdFromItem(&item)
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	itemStruct := Item{}
	err = json.Unmarshal([]byte(item), &itemStruct)
	if err != nil {
		fmt.Println(err)
		return
	}
	var content []byte
	if len(fileContent) == 0 {
		content, err = json.Marshal([]Item{itemStruct})
	} else {
		var items []Item
		err = json.Unmarshal(fileContent, &items)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, item := range items {
			if item.Id == id {
				msg := []byte(fmt.Sprintf("Item with id %s already exists", id))
				writer.Write(msg)
				return
			}
		}
		items = append(items, itemStruct)
		content, err = json.Marshal(&items)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(content)
}

func list(filename string, writer io.Writer) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(content) == 0 {
		return
	}
	writer.Write(content)
}

func findById(filename, id string, writer io.Writer) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var items []Item
	err = json.Unmarshal(fileContent, &items)
	if err != nil {
		fmt.Println(err)
		return
	}
	var content []byte
	for _, item := range items {
		if item.Id == id {
			content, err = json.Marshal(item)
			if err != nil {
				fmt.Println(err)
				return
			}
			writer.Write(content)
			return
		}
	}
	writer.Write(content)
}

func remove(filename, id string, writer io.Writer) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var items []Item
	err = json.Unmarshal(fileContent, &items)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, item := range items {
		if item.Id == id {
			items = append(items[:i], items[i+1:]...)
			content, err := json.Marshal(items)
			if err != nil {
				fmt.Println(err)
				return
			}
			file.Truncate(0)
			file.Seek(0, 0)
			file.Write(content)
			return
		}
	}
	msg := []byte(fmt.Sprintf("Item with id %s not found", id))
	writer.Write(msg)
}

func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]
	item := args["item"]
	fileName := args["fileName"]
	id := args["id"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}
	switch operation {
	case "add":
		if item == "" {
			return errors.New("-item flag has to be specified")
		}
		add(item, fileName, writer)
	case "list":
		list(fileName, writer)
	case "findById":
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		findById(fileName, id, writer)
	case "remove":
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		remove(fileName, id, writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	return nil
}

func parseIdFromItem(item *string) (string, error) {
	itemStruct := Item{}
	err := json.Unmarshal([]byte(*item), &itemStruct)
	return itemStruct.Id, err
}

func main() {
	item := flag.String("item", "", `Usage example: '{"id": "1", "email": "email@test.com", "age": 23}'`)
	operation := flag.String("operation", "", `Possible values: 'add', 'list', 'findById', 'remove'`)
	fileName := flag.String("fileName", "", `Usage example: 'users.json'`)
	id := flag.String("id", "", "An integer value > 0")
	flag.Parse()
	args := map[string]string{
		"item":      *item,
		"operation": *operation,
		"fileName":  *fileName,
		"id":        *id,
	}
	err := Perform(args, os.Stdout)
	if err != nil {
		panic(err)
	}
}
