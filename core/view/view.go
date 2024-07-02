package view

import (
	"app/core/data"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"io"
	"log"
	"os"
)

type Data data.Map

var (
	Path = os.Getenv("VIEW_PATH")
	Ext  = os.Getenv("VIEW_EXT")
)

// LoadView load view content from template.
func LoadView(template string, d Data) string {
	file, err := pongo2.FromFile(fmt.Sprintf("%s/%s.%s", Path, template, Ext))
	if err != nil {
		log.Fatal(err)
	}
	ctx := pongo2.Context(d)
	res, err := file.Execute(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

// LoadViewWriter load view content into writer.
func LoadViewWriter(template string, d Data, writer io.Writer) error {
	var templateFile = pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s/%s.%s", Path, template, Ext)))
	ctx := pongo2.Context(d)

	return templateFile.ExecuteWriter(ctx, writer)
}
