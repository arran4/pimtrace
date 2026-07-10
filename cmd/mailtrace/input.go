package main

import (
	"fmt"
	"io"
	_ "github.com/emersion/go-message/charset"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/maildata"
	"pimtrace/fsys"
)

func InputHandler(fs fsys.FS, inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
	var out io.Writer = os.Stdout
	for _, op := range ops {
		if w, ok := op.(io.Writer); ok && w != nil {
			out = w
			break
		}
	}

	var mails []*maildata.MailWithSource
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
			nm, err := dataformats.ReadFile(fs, inputType, inputFile, maildata.ReadMailStream)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "mboxgz":
		ops = append(ops, dataformats.Gzip)
		fallthrough
	case "mbox":
		switch inputFile {
		case "-":
			nm, err := maildata.ReadMBoxStream(os.Stdin, inputType, inputFile, ops...)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := dataformats.ReadFile(fs, inputType, inputFile, maildata.ReadMBoxStream, ops...)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "mboxtargz":
		ops = append(ops, dataformats.Gzip)
		fallthrough
	case "mboxtar":
		switch inputFile {
		case "-":
			nm, err := dataformats.ReadTarStream(os.Stdin, inputType, inputFile, maildata.ReadMBoxStream, []string{"*.mbox"}, ops...)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := dataformats.ReadTarFile(fs, inputType, inputFile, maildata.ReadMBoxStream, []string{"*.mbox"}, ops...)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "list":
		PrintInputHelp(out)
	default:
		return nil, fmt.Errorf("please specify an -input-type. got %s", inputType)
	}
	return maildata.Data(mails), nil
}

func PrintInputHelp(w io.Writer) {
	_, _ = fmt.Fprintln(w, "input-types available: ")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "mailfile", "A single mail file")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "mbox", "Mbox file")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "mboxgz", "Gzipped Mbox file")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "mboxtargz", "Gzipped Tarred collection of Mbox file")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "list", "This help text")
	_, _ = fmt.Fprintln(w)
}
