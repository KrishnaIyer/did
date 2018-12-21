package did // import "go.htdvisser.nl/did"

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/spf13/pflag"
	"go.htdvisser.nl/did/day"
	"go.htdvisser.nl/did/db"
)

var (
	flags  = pflag.NewFlagSet(filepath.Base(os.Args[0]), pflag.ContinueOnError)
	dbFile = flags.String("db", "$HOME/.did.db", "db file")
	format = flags.String("format", `{{.GetTime.Local.Format "15:04"}}: {{.GetMessage}}`, "format template")
)

func init() {
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintf(os.Stderr, "  %s [%s]\n", filepath.Base(os.Args[0]), strings.Join(day.Indicators(), "|"))
		fmt.Fprintf(os.Stderr, "  %s [description]\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Flags:")
		flags.PrintDefaults()
	}
	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func Main() (err error) {
	args := flags.Args()
	if len(args) < 1 {
		args = []string{"today"}
	}
	db, err := db.Open(os.ExpandEnv(*dbFile))
	if err != nil {
		return err
	}
	defer func() {
		closeErr := db.Close()
		if err == nil {
			err = closeErr
		}
	}()
	selected, ok := day.GetMidnight(args[0])
	if ok {
		return history(db, selected)
	}
	return addRecord(db, args...)
}

var jsonMarshaler = &jsonpb.Marshaler{OrigName: true}

func history(didDB *db.DB, day time.Time) error {
	var render func(*db.Record) error
	switch *format {
	case "json": // JSON stream.
		render = func(record *db.Record) error {
			return jsonMarshaler.Marshal(os.Stdout, record)
		}
	default: // Render each record with a template.
		template, err := template.New("format").Parse(*format)
		if err != nil {
			return err
		}
		render = func(record *db.Record) error {
			return template.Execute(os.Stdout, &record)
		}
	}
	history, err := didDB.History(day)
	if err != nil {
		return err
	}
	for _, record := range history {
		if err := render(record); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func addRecord(didDB *db.DB, args ...string) error {
	timestamp, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return err
	}
	return didDB.AddRecord(&db.Record{
		Timestamp: timestamp,
		Message:   strings.Join(args, " "),
	})
}
