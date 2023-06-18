package main

import (
	"flag"
	"fmt"
	"os"

	cli "github.com/AnkitKumar117/fileStore/pkg/cli"
)

const A_ADD, A_LIST, A_REMOVE, A_UPDATE, A_WC, A_FREQ_WORDS = "add", "ls", "rm", "update", "wc", "freq-words"

func main() {
	/*
		store add file1.txt file2.txt, --> working fine
		store ls, --> DONE
		store rm file.txt, --> DONE
		store update file.txt, --> done
		store wc, --> DONE
		store freq-words [--limit|-n 10] [--order=dsc|asc]   |--> DONE -- check order once
	*/
	wcCommand := flag.NewFlagSet("wc", flag.ExitOnError)
	lsCommand := flag.NewFlagSet("ls", flag.ExitOnError)
	rmCommand := flag.NewFlagSet("rm", flag.ExitOnError)
	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	fwCommand := flag.NewFlagSet("freq-words", flag.ExitOnError)

	var limit int
	var order string
	fwCommand.IntVar(&limit, "limit", 10, "Default limit is 10")
	fwCommand.IntVar(&limit, "n", 10, "Default limit is 10")
	fwCommand.StringVar(&order, "order", "desc", "Default limit is desc")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Action is required..")
		os.Exit(1)
	}
	fmt.Printf("File store assessment for redhat\n")
	action := os.Args[1]
	fmt.Println("Action : ", action)
	switch action {
	case A_ADD:
		addCommand.Parse(os.Args[2:])
		cli.AddFiles(addCommand.Args())
	case A_LIST:
		lsCommand.Parse(os.Args[2:])
		if len(lsCommand.Args()) > 0 {
			fmt.Println("There is no sub-command/arguments of ls. Try -- ./store ls ")
			os.Exit(1)
		}
		cli.AllList()
	case A_REMOVE:
		rmCommand.Parse(os.Args[2:])
		cli.RemoveFile(rmCommand.Args())
	case A_UPDATE:
		updateCommand.Parse(os.Args[2:])
		cli.UpdateFile(updateCommand.Args())
	case A_WC:
		wcCommand.Parse(os.Args[2:])
		if len(wcCommand.Args()) > 0 {
			fmt.Println("There is no sub-command/arguments of wc. Try -- ./store wc ")
			os.Exit(1)
		}
		cli.WordCount()
	case A_FREQ_WORDS:
		fwCommand.Parse(os.Args[2:])
		cli.FreqWords(limit, order)
	default:
		fmt.Println(" Valid actions are :", A_ADD, A_LIST, A_REMOVE, A_UPDATE, A_WC, A_FREQ_WORDS)
	}

}
