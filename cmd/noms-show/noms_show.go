// Copyright 2016 The Noms Authors. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/outputpager"
)

var (
	showHelp = flag.Bool("help", false, "show help text")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Shows a serialization of a Noms object\n")
		fmt.Fprintln(os.Stderr, "Usage: noms show <object>\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSee \"Spelling Objects\" at https://github.com/attic-labs/noms/blob/master/doc/spelling.md for details on the object argument.\n\n")
	}

	flag.Parse()
	if *showHelp {
		flag.Usage()
		return
	}

	if len(flag.Args()) != 1 {
		d.CheckErrorNoUsage(errors.New("expected exactly one argument"))
	}

	database, value, err := spec.GetPath(flag.Arg(0))
	d.CheckErrorNoUsage(err)

	if value == nil {
		fmt.Fprintf(os.Stderr, "Object not found: %s\n", flag.Arg(0))
		return
	}

	waitChan := outputpager.PageOutput(!*outputpager.NoPager)

	w := bufio.NewWriter(os.Stdout)
	types.WriteEncodedValueWithTags(w, value)
	fmt.Fprintf(w, "\n")
	w.Flush()
	database.Close()

	if waitChan != nil {
		os.Stdout.Close()
		<-waitChan
	}
}
