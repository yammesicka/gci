package gci

type Cases struct {
	name, config, in, out string
}

var commonConfig = `sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
`

var testCases = []Cases{
	{
		"already-good",
		commonConfig,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"blank-format",
		commonConfig,
		`package main
import (
	"fmt"

  // comment
	g  "github.com/golang"    // comment

	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	// comment
	g "github.com/golang" // comment

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"cgo-block",
		commonConfig,
		`package main

import (
	/*
	#include "types.h"
	*/
	"C"
)
`,
		`package main

import (
	/*
	#include "types.h"
	*/
	"C"
)
`,
	},
	{
		"cgo-block-after-import",
		commonConfig,
		`package main

import (
	"fmt"

	"github.com/yammesicka/gci"
	g "github.com/golang"
)

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"
`,
		`package main

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"cgo-block-before-import",
		commonConfig,
		`package main

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"

import (
	"fmt"

	"github.com/yammesicka/gci"

	g "github.com/golang"
)
`,
		`package main

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"cgo-block-mixed",
		commonConfig,
		`package main

import (
	/* #include "types.h"
	*/"C"
)
`,
		`package main

import (
	/* #include "types.h"
	*/"C"
)
`,
	},
	{
		"cgo-block-mixed-with-content",
		commonConfig,
		`package main

import (
	/* #include "types.h"
	#include "other.h" */"C"
)
`,
		`package main

import (
	/* #include "types.h"
	#include "other.h" */"C"
)
`,
	},
	{
		"cgo-block-prefix",
		commonConfig,
		`package main

import (
	/* #include "types.h" */ "C"
)
`,
		`package main

import (
	/* #include "types.h" */ "C"
)
`,
	},
	{
		"cgo-block-single-line",
		commonConfig,
		`package main

import (
	/* #include "types.h" */
	"C"
)
`,
		`package main

import (
	/* #include "types.h" */
	"C"
)
`,
	},
	{
		"cgo-line",
		commonConfig,
		`package main

import (
	// #include "types.h"
	"C"
)
`,
		`package main

import (
	// #include "types.h"
	"C"
)
`,
	},
	{
		"cgo-multiline",
		commonConfig,
		`package main

import (
	// #include "types.h"
	// #include "other.h"
	"C"
)
`,
		`package main

import (
	// #include "types.h"
	// #include "other.h"
	"C"
)
`,
	},
	{
		"cgo-single",
		commonConfig,
		`package main

import (
	"fmt"

	"github.com/yammesicka/gci"
)

import "C"

import "github.com/golang"

import (
  "github.com/yammesicka/gci"
)
`,
		`package main

import "C"

import (
	"fmt"

	"github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"comment",
		commonConfig,
		`package main
import (
	//Do not forget to run Gci
	"fmt"
)
`,
		`package main
import (
	//Do not forget to run Gci
	"fmt"
)
`,
	},
	{
		"comment-before-import",
		commonConfig,
		`package main

// comment
import (
	"fmt"
	"os"

	"github.com/yammesicka/gci"
)
`,
		`package main

// comment
import (
	"fmt"
	"os"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"comment-in-the-tail",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)

type test int

// test
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)

type test int

// test
`,
	},
	{
		"comment-top",
		commonConfig,
		`package main

import (
	"os" // https://pkg.go.dev/os
	// https://pkg.go.dev/fmt
	"fmt"
)
`,
		`package main

import (
	// https://pkg.go.dev/fmt
	"fmt"
	"os" // https://pkg.go.dev/os
)
`,
	},
	{
		"comment-without-whitespace",
		commonConfig,
		`package proc

import (
	"context"// no separating whitespace here //nolint:confusion
)
`,
		`package proc

import (
	"context"// no separating whitespace here //nolint:confusion
)
`,
	},
	{
		"comment-with-slashslash",
		commonConfig,
		`package main

import (
	"fmt" // https://pkg.go.dev/fmt
)
`,
		`package main

import (
	"fmt" // https://pkg.go.dev/fmt
)
`,
	},
	{
		"custom-order",
		`customOrder: true
sections:
  - Prefix(github.com/yammesicka)
  - Default
  - Standard
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/a"
)
`,
		`package main

import (
	"github.com/yammesicka/a"

	g "github.com/golang"

	"fmt"
)
`,
	},
	{
		"default-order",
		`sections:
  - Standard
  - Prefix(github.com/yammesicka)
  - Default
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/a"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/a"
)
`,
	},
	{
		"dot-and-blank",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
  - Blank
  - Dot
`,
		`package main

import (
	"fmt"

	g "github.com/golang"
	. "github.com/golang/dot"
	_ "github.com/golang/blank"

	"github.com/yammesicka/a"
	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
	. "github.com/yammesicka/gci/dot"
	_ "github.com/yammesicka/gci/blank"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/a"
	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"

	_ "github.com/yammesicka/gci/blank"
	_ "github.com/golang/blank"

	. "github.com/yammesicka/gci/dot"
	. "github.com/golang/dot"
)
`,
	},
	{
		"duplicate-imports",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	a "github.com/yammesicka/gci"
	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
	a "github.com/yammesicka/gci"
)
`,
	},
	{
		"grouped-multiple-custom",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka,gitlab.com/yammesicka,yammesicka)
`,
		`package main

import (
	"yammesicka/lib1"
	"fmt"
	"github.com/yammesicka/gci"
	"gitlab.com/yammesicka/gci"
	g "github.com/golang"
	"github.com/yammesicka/gci/subtest"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"yammesicka/lib1"
	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
	"gitlab.com/yammesicka/gci"
)
`,
	},
	{
		"leading-comment",
		commonConfig,
		`package main

import (
	// foo
	"fmt"
)
`,
		`package main

import (
	// foo
	"fmt"
)
`,
	},
	{
		"linebreak",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
`,
		`package main

import (
	g "github.com/golang"

	"fmt"

	"github.com/yammesicka/gci"

)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"linebreak-no-custom",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
`,
		`package main

import (
	g "github.com/golang"

	"fmt"

)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"
)
`,
	},
	{
		"mismatch-section",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
  - Prefix(github.com/yammesicka/gci)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"multiple-custom",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
  - Prefix(github.com/yammesicka/gci)
  - Prefix(github.com/yammesicka/gci/subtest)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/a"
	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/a"

	"github.com/yammesicka/gci"

	"github.com/yammesicka/gci/subtest"
)
`,
	},
	{
		"multiple-imports",
		commonConfig,
		`package main

import "fmt"

import "context"

import (
	"os"

	"github.com/yammesicka/test"
)

import "math"


// main
func main() {
}
`,
		`package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/yammesicka/test"
)

// main
func main() {
}
`,
	},
	{
		"multiple-line-comment",
		commonConfig,
		`package proc

import (
	"context" // in-line comment
	"fmt"
	"os"

	//nolint:depguard // A multi-line comment explaining why in
	// this one case it's OK to use os/exec even though depguard
	// is configured to force us to use dlib/exec instead.
	"os/exec"

	"golang.org/x/sys/unix"
	"github.com/local/dlib/dexec"
)
`,
		`package proc

import (
	"context" // in-line comment
	"fmt"
	"os"
	//nolint:depguard // A multi-line comment explaining why in
	// this one case it's OK to use os/exec even though depguard
	// is configured to force us to use dlib/exec instead.
	"os/exec"

	"github.com/local/dlib/dexec"
	"golang.org/x/sys/unix"
)
`,
	},
	{
		"nochar-after-import",
		commonConfig,
		`package main

import (
	"fmt"
)
`,
		`package main

import (
	"fmt"
)
`,
	},
	{
		"no-format",
		commonConfig,
		`package main

import(
"fmt"

g "github.com/golang"

"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"nolint",
		commonConfig,
		`package main

import (
	"fmt"

	"github.com/forbidden/pkg" //nolint:depguard

	_ "github.com/yammesicka/gci" //nolint:depguard
)
`,
		`package main

import (
	"fmt"

	"github.com/forbidden/pkg" //nolint:depguard

	_ "github.com/yammesicka/gci" //nolint:depguard
)
`,
	},
	{
		"number-in-alias",
		commonConfig,
		`package main

import (
	"fmt"

	go_V1 "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	go_V1 "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"one-import",
		commonConfig,
		`package main
import (
	"fmt"
)

func main() {
}
`,
		`package main
import (
	"fmt"
)

func main() {
}
`,
	},
	{
		"one-import-one-line",
		commonConfig,
		`package main

import "fmt"

func main() {
}
`,
		`package main

import "fmt"

func main() {
}
`,
	},
	{
		"one-line-import-after-import",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka)
`,
		`package main

import (
	"fmt"
	"os"

	"github.com/yammesicka/test"
)

import "context"
`,
		`package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yammesicka/test"
)
`,
	},
	{
		"same-prefix-custom",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka/gci)
  - Prefix(github.com/yammesicka/gci/subtest)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"

	"github.com/yammesicka/gci/subtest"
)
`,
	},
	{
		"same-prefix-custom",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka/gci)
  - Prefix(github.com/yammesicka/gci/subtest)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"

	"github.com/yammesicka/gci/subtest"
)
`,
	},
	{
		"same-prefix-custom",
		`sections:
  - Standard
  - Default
  - Prefix(github.com/yammesicka/gci)
  - Prefix(github.com/yammesicka/gci/subtest)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"

	"github.com/yammesicka/gci/subtest"
)
`,
	},
	{
		"simple-case",
		commonConfig,
		`package main

import (
	"golang.org/x/tools"

	"fmt"

	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	"golang.org/x/tools"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"whitespace-test",
		commonConfig,
		`package main

import (
	"fmt"
	"github.com/golang" // golang
	alias "github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	"github.com/golang" // golang

	alias "github.com/yammesicka/gci"
)
`,
	},
	{
		"with-above-comment-and-alias",
		commonConfig,
		`package main

import (
	"fmt"
	// golang
	_ "github.com/golang"
	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	// golang
	_ "github.com/golang"

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"with-comment-and-alias",
		commonConfig,
		`package main

import (
	"fmt"
	_ "github.com/golang" // golang
	"github.com/yammesicka/gci"
)
`,
		`package main

import (
	"fmt"

	_ "github.com/golang" // golang

	"github.com/yammesicka/gci"
)
`,
	},
	{
		"blank-in-config",
		`sections:
  - Standard
  - Default
  - Prefix(  github.com/yammesicka/gci,   github.com/yammesicka/gci/subtest  )
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
)
`,
		`package main

import (
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
)
`,
	},
	{
		"alias",
		`sections:
  - Standard
  - Default
  - Alias
`,
		`package main

import (
	testing "github.com/yammesicka/test"
	"fmt"

	g "github.com/golang"

	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"
)
`,
		`package main

import (
	"fmt"

	"github.com/yammesicka/gci"
	"github.com/yammesicka/gci/subtest"

	testing "github.com/yammesicka/test"
	g "github.com/golang"
)
`,
	},
	{
		"no-trailing-newline",
		`sections:
  - Standard
`,
		`package main

import (
	"net"
	"fmt"
)`,
		`package main

import (
	"fmt"
	"net"
)
`,
	},
}
