// Copyright © 2011-12 Qtrac Ltd.
// 
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "compress/gzip"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "strconv"
    "time"
)

const (
    fileType             = "INVOICES"   // Used by text formats
    magicNumber          = 0x125D       // Used by binary formats
    fileVersion          = 101          // Used by all formats
    dateFormat           = "2006-01-02" // This date must always be used
    nanosecondsToSeconds = 1e9
)

type Invoice struct {
    Id         		int
    CustomerId 		int
    DepartmentId	string
    Raised     		time.Time
    Due        		time.Time
    Paid       		bool
    Note       		string
    Items      		[]*Item
}

type Item struct {
    Id       	string
    Price    	float64
    Quantity 	int
    TaxBand		int
    Note     	string
}

type InvoicesMarshaler interface {
    MarshalInvoices(writer io.Writer, invoices []*Invoice) error
}

type InvoicesUnmarshaler interface {
    UnmarshalInvoices(reader io.Reader) ([]*Invoice, error)
}

func main() {
    log.SetFlags(0)
    report := false
    args := os.Args[1:]
    if len(args) > 0 && (args[0] == "-t" || args[0] == "--time") {
        report = true
        args = args[1:]
    }
    if len(args) != 2 || args[0] == "-h" || args[0] == "--help" {
        log.Fatalf("usage: %s [-t|--time] infile.ext outfile.ext\n"+
            ".ext may be any of .gob, .inv, .jsn, .json, .txt, "+
            "or .xml, optionally gzipped (e.g., .gob.gz)\n",
            filepath.Base(os.Args[0]))
    }
    inFilename, outFilename := args[0], args[1]
    if inFilename == outFilename {
        log.Fatalln("won't overwrite a file with itself")
    }

    start := time.Now()
    invoices, err := readInvoiceFile(inFilename)
    if err == nil && report {
        duration := time.Now().Sub(start)
        fmt.Printf("Read  %s in %.3f seconds\n", inFilename,
            float64(duration)/nanosecondsToSeconds)
    }
    if err != nil {
        log.Fatalln("Failed to read:", err)
    }
    
    if invoices, err = updateInvoices(invoices); err != nil {
    	log.Fatalln("Failed to update:", err)
    }
    
    start = time.Now()
    err = writeInvoiceFile(outFilename, invoices)
    if err == nil && report {
        duration := time.Now().Sub(start)
        fmt.Printf("Wrote %s in %.3f seconds\n", outFilename,
            float64(duration)/nanosecondsToSeconds)
    }
    if err != nil {
        log.Fatalln("Failed to write:", err)
    }
}

func updateInvoices(invoices []*Invoice) ([]*Invoice, error) {
	updated := make([]*Invoice, 0, len(invoices))
	var err error
	for _, invoice := range invoices {
		if invoice, err = updateInvoiceData(invoice); err != nil {
			return nil, err
		}
		updated = append(updated, invoice)
	}
	return updated, nil
}

func getDepartId(id int) string {
	var dptId string
	switch {
		case id < 3000:
			dptId = "GEN"
		case id < 4000:
			dptId = "MKT"
		case id < 5000:
			dptId = "COM"
		case id < 6000:
			dptId = "EXP"
		case id < 7000:
			dptId = "INP"
		case id < 8000:
			dptId = "TZZ"
		case id < 9000:
			dptId = "V20"
		default:
			dptId = "X15"
	}
	return dptId
}

func updateInvoiceData(invoice *Invoice) (*Invoice, error) {
	invoice.DepartmentId = getDepartId(invoice.Id)
	var err error
	for _, item := range invoice.Items {
		if item, err = updateItem(item); err != nil {
			return nil, err
		}
	}
	
	return invoice, nil
}

func updateItem(item *Item) (*Item, error) {
	txbd, err := strconv.Atoi(string(item.Id[2]))
	if err != nil {
		return nil, err
	}
	item.TaxBand = txbd
	return item, nil
}

func readInvoiceFile(filename string) ([]*Invoice, error) {
    file, closer, err := openInvoiceFile(filename)
    if closer != nil {
        defer closer()
    }
    if err != nil {
        return nil, err
    }
    return readInvoices(file, suffixOf(filename))
}

func openInvoiceFile(filename string) (io.ReadCloser, func(), error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, err
    }
    closer := func() { file.Close() }
    var reader io.ReadCloser = file
    var decompressor *gzip.Reader
    if strings.HasSuffix(filename, ".gz") {
        if decompressor, err = gzip.NewReader(file); err != nil {
            return file, closer, err
        }
        closer = func() { decompressor.Close(); file.Close() }
        reader = decompressor
    }
    return reader, closer, nil
}

func readInvoices(reader io.Reader, suffix string) ([]*Invoice, error) {
    var unmarshaler InvoicesUnmarshaler
    switch suffix {
    case ".gob":
        unmarshaler = GobMarshaler{}
    case ".inv":
        unmarshaler = InvMarshaler{}
    case ".jsn", ".json":
        unmarshaler = JSONMarshaler{}
    case ".txt":
        unmarshaler = TxtMarshaler{}
    case ".xml":
        unmarshaler = XMLMarshaler{}
    }
    if unmarshaler != nil {
        return unmarshaler.UnmarshalInvoices(reader)
    }
    return nil, fmt.Errorf("unrecognized input suffix: %s", suffix)
}

func writeInvoiceFile(filename string, invoices []*Invoice) error {
    file, closer, err := createInvoiceFile(filename)
    if closer != nil {
        defer closer()
    }
    if err != nil {
        return err
    }
    return writeInvoices(file, suffixOf(filename), invoices)
}

func createInvoiceFile(filename string) (io.WriteCloser, func(), error) {
    file, err := os.Create(filename)
    if err != nil {
        return nil, nil, err
    }
    closer := func() { file.Close() }
    var writer io.WriteCloser = file
    var compressor *gzip.Writer
    if strings.HasSuffix(filename, ".gz") {
        compressor = gzip.NewWriter(file)
        closer = func() { compressor.Close(); file.Close() }
        writer = compressor
    }
    return writer, closer, nil
}

func writeInvoices(writer io.Writer, suffix string,
    invoices []*Invoice) error {
    var marshaler InvoicesMarshaler
    switch suffix {
    case ".gob":
        marshaler = GobMarshaler{}
    case ".inv":
        marshaler = InvMarshaler{}
    case ".jsn", ".json":
        marshaler = JSONMarshaler{}
    case ".txt":
        marshaler = TxtMarshaler{}
    case ".xml":
        marshaler = XMLMarshaler{}
    }
    if marshaler != nil {
        return marshaler.MarshalInvoices(writer, invoices)
    }
    return errors.New("unrecognized output suffix")
}

func suffixOf(filename string) string {
    suffix := filepath.Ext(filename)
    if suffix == ".gz" {
        suffix = filepath.Ext(filename[:len(filename)-3])
    }
    return suffix
}
