package main

import (
	"fmt"
	_ "github.com/emersion/go-message/charset"
	"os"
	"pimtrace"
	"pimtrace/dataformats/maildata"
)

func InputHandler(inputType string, inputFile string) (pimtrace.Data, error) {
	mails := []*maildata.MailWithSource{}
	switch inputType {
	case "mailfile":
		switch inputFile {
		case "-":
			nm, err := maildata.ReadMailStream(os.Stdin, inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := maildata.ReadMailFile(inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "mbox":
		switch inputFile {
		case "-":
			nm, err := maildata.ReadMBoxStream(os.Stdin, inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := maildata.ReadMBoxFile(inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "list":
		PrintInputHelp()
	default:
		return nil, fmt.Errorf("please specify an -input-type. got %s", inputType)
	}
	return maildata.Data(mails), nil
}

func PrintInputHelp() {
	fmt.Println("`input-type`s available: ")
	fmt.Printf(" %-30s %s\n", "mailfile", "A single mail file")
	fmt.Printf(" %-30s %s\n", "mbox", "Mbox file")
	fmt.Printf(" %-30s %s\n", "list", "This help text")
	fmt.Println()
}
