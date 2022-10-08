package main

import (
	"fmt"
	"os"
)

func OutputHandler(mails []*MailWithSource) error {
	switch *inputType {
	case "mailfile":
		switch *inputFile {
		case "-":
			nm, err := WriteMailStream(os.Stdin, *inputType, *inputFile)
			if err != nil {
				return err
			}
			mails = append(mails, nm...)
		default:
			nm, err := WriteMailFile(*inputType, *inputFile)
			if err != nil {
				return err
			}
			mails = append(mails, nm...)
		}
	case "mbox":
		switch *inputFile {
		case "-":
			nm, err := WriteMBoxStream(os.Stdin, *inputType, *inputFile)
			if err != nil {
				return err
			}
			mails = append(mails, nm...)
		default:
			nm, err := WriteMBoxFile(*inputType, *inputFile)
			if err != nil {
				return err
			}
			mails = append(mails, nm...)
		}
	case "list":
		fmt.Println("`--input-type`s: ")
		fmt.Printf(" =%-20s - %s\n", "mailfile", "A single mail file")
		fmt.Printf(" =%-20s - %s\n", "mbox", "Mbox file")
		fmt.Printf(" =%-20s - %s\n", "list", "This help text")
		fmt.Println()
	default:
		fmt.Println("Please specify a -input-type")
		fmt.Println()
	}
	return nil
}
