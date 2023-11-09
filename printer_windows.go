/* printer_windows.go */

package utopia

import (
	"log"

	prt "github.com/alexbrainman/printer"
)

// doctype "RAW", "XPS_PASS", "TEXT"

func GoPrinter() (hObj *GoPrinterObj) {
	name, err := prt.Default() // returns name of Default Printer as string
	if err != nil {
	    log.Println("No default printer is available...", err)
	}
	device, err := prt.Open(name) // Opens the named printer and returns a *Printer
	if err != nil {
	    log.Println("Failed to access printer...", err)
	}
	return &GoPrinterObj{name, device, false}
}

type GoPrinterObj struct {
	name 	string	// default
	dev 	*prt.Printer
	inpage 	bool
}

func (ob *GoPrinterObj) CloseDocument() (err error) {
	err = ob.dev.Close()
	if err != nil {
		log.Println("Failed to close document...", err)
	}
	return err
}

func (ob *GoPrinterObj) EndPage() (err error) {
	err = ob.dev.EndPage()
	if err != nil {
		log.Println("Failed to end page...", err)
	}
	return err
}

func (ob *GoPrinterObj) NewPage() (err error) {
	err = ob.dev.StartPage()
	if err != nil {
		log.Println("Failed to create new page...", err)
	}
	return err
}

func (ob *GoPrinterObj) OpenDocument(docname string, doctype string) (err error) {
	err = ob.dev.StartDocument(docname, doctype)
	if err != nil {
		log.Println("Failed to open new document...", err)
	}
	return err
}

/*func (ob *GoPrinterObj) AvailablePrinters() (printers []string, err error) {
	printers, err = ReadNames()
	return
}*/

func (ob *GoPrinterObj) Write(text []byte) (numbytes int, err error) {
	numbytes, err = ob.dev.Write(text)
	if err != nil {
		log.Println("Failed to write to page...", err)
	}
	return numbytes, err
}

