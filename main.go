package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func InitFlag() (string, string, string, string) {
	id := flag.String("id", "", "DocID")
	toggle := flag.String("toggle", "", "On or Off")
	cronjob := flag.String("cronjob", "", "Cronjob")
	op := flag.String("op", "", "CRUD")
	flag.Parse()
	return *id, *toggle, *cronjob, *op
}

func GetCronFile() {
	cmd := exec.Command("bash", "-c", "crontab -l > file")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func SetCronFile() {
	cmd := exec.Command("bash", "-c", "crontab file")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func FileToVar(n string, v *[]string) {
	file, err := os.Open(n)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		*v = append(*v, scanner.Text())
	}
	ferr := file.Close()
	if ferr != nil {
		log.Fatal(ferr)
	}
}

func SliceToString(f *[]string) string {
	var temp string
	for _, s := range *f {
		temp += s + "\n"
	}
	return temp
}

func WriteToFile(name string, data string) {
	err := os.WriteFile(name, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func CUFile(l *[]string, id string, toggle string, cronjob string) bool {
	var check bool
	for i, data := range *l {
		if strings.Contains(data, fmt.Sprintf("#%v:%v", id, toggle)) {
			check = true
			(*l)[i+1] = cronjob
			break
		}
		check = false
	}
	if check {
		fmt.Println("File updated")
	} else {
		fmt.Println("Cronjob added")
		*l = append(*l, fmt.Sprintf("#%v:%v", id, toggle), cronjob)
	}
	return check
}

func DFile(l *[]string, id string, toggle string) bool {
	var check bool
	for i, data := range *l {
		if strings.Contains(data, fmt.Sprintf("#%v:%v", id, toggle)) {
			check = true
			*l = append((*l)[:i], (*l)[i+2:]...)
			break
		}
		check = false
	}
	if check {
		fmt.Println("File updated")
	} else {
		fmt.Println(fmt.Sprintf("Couldnt find #%v:%v", id, toggle))
	}
	return check
}

var line []string

func main() {

	id, toggle, cronjob, op := InitFlag()

	GetCronFile()

	FileToVar("file", &line)

	if op == "create" || op == "update" {
		result := CUFile(&line, id, toggle, cronjob)
		fmt.Println(result)
	}
	if op == "delete" {
		result := DFile(&line, id, toggle)
		fmt.Println(result)
	}

	fmt.Println(cronjob)

	WriteToFile("file", SliceToString(&line))

	SetCronFile()
}
