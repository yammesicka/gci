package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

const C = "\"C\""

type GciImports struct {
	// original index of import group, include doc, name, path and comment
	Start, End int
	Name, Path string
}
type ImportList []*GciImports

func (l ImportList) Len() int {
	return len(l)
}

func (l ImportList) Less(i, j int) bool {
	if strings.Compare(l[i].Path, l[j].Path) == 0 {
		return strings.Compare(l[i].Name, l[j].Name) < 0
	}

	return strings.Compare(l[i].Path, l[j].Path) < 0
}

func (l ImportList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

/*
 * AST considers a import block as below:
 * ```
 * Doc
 * Name Path Comment
 * ```
 * An example is like below:
 * ```
 * // test
 * test "fmt" // test
 * ```
 * getImports return a import block with name, start and end index
 */
func getImports(imp *ast.ImportSpec) (start, end int, name string) {
	if imp.Doc != nil {
		// doc poc need minus one to get the first index of comment
		start = int(imp.Doc.Pos()) - 1
	} else {
		if imp.Name != nil {
			// name pos need minus one too
			start = int(imp.Name.Pos()) - 1
		} else {
			// path pos start without quote, need minus one for it
			start = int(imp.Path.Pos()) - 1
		}
	}

	if imp.Name != nil {
		name = imp.Name.Name
	}

	if imp.Comment != nil {
		end = int(imp.Comment.End())
	} else {
		end = int(imp.Path.End())
	}
	return
}

// importDeclRange tracks the byte range of a single import declaration (GenDecl).
type importDeclRange struct {
	start, end int
}

func ParseFile(src []byte, filename string) (ImportList, int, int, int, int, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	if len(f.Imports) == 0 {
		return nil, 0, 0, 0, 0, NoImportError{}
	}

	var (
		// headEnd means the start of import block
		headEnd int
		// tailStart means the end + 1 of import block
		tailStart int
		// cStart means the start of C import block
		cStart int
		// cEnd means the end of C import block
		cEnd int
		data ImportList
		// importDecls tracks byte ranges of each non-C import declaration,
		// used to detect standalone comments between separate import declarations.
		importDecls []importDeclRange
	)

	for index, decl := range f.Decls {
		switch decl.(type) {
		// skip BadDecl and FuncDecl
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			if genDecl.Tok == token.IMPORT {
				// there are two cases, both end with linebreak:
				// 1.
				// import (
				//	 "xxxx"
				// )
				// 2.
				// import "xxx"
				if headEnd == 0 {
					headEnd = int(decl.Pos()) - 1
				}
				tailStart = int(decl.End())
				if tailStart > len(src) {
					tailStart = len(src)
				}

				isCImport := false
				for _, spec := range genDecl.Specs {
					imp := spec.(*ast.ImportSpec)
					// there are only one C import block
					// ensure C import block is the first import block
					if imp.Path.Value == C {
						/*
							common case:

							// #include <png.h>
							import "C"

							notice that decl.Pos() == genDecl.Pos() > genDecl.Doc.Pos()
						*/
						if genDecl.Doc != nil {
							cStart = int(genDecl.Doc.Pos()) - 1
							// if C import block is the first, update headEnd
							if index == 0 {
								headEnd = cStart
							}
						} else {
							/*
								special case:

								import "C"
							*/
							cStart = int(decl.Pos()) - 1
						}

						cEnd = int(decl.End())
						isCImport = true

						continue
					}

					start, end, name := getImports(imp)

					data = append(data, &GciImports{
						Start: start,
						End:   end,
						Name:  name,
						Path:  strings.Trim(imp.Path.Value, `"`),
					})
				}

				// Track non-C import declaration ranges for standalone comment detection.
				if !isCImport {
					declStart := int(decl.Pos()) - 1
					declEnd := int(decl.End())
					if declEnd > len(src) {
						declEnd = len(src)
					}
					importDecls = append(importDecls, importDeclRange{start: declStart, end: declEnd})
				}
			}
		}
	}

	// If there are multiple separate non-C import declarations, check for
	// standalone comments between them. If found, skip reformatting to
	// preserve the original structure and avoid dropping the comments.
	if len(importDecls) > 1 {
		if hasCommentBetweenDecls(f.Comments, importDecls, cStart, cEnd) {
			return nil, 0, 0, 0, 0, CommentBetweenImportsError{}
		}
	}

	// Attach standalone comments within the import block to adjacent imports.
	// This handles comments separated by blank lines from imports, which the
	// Go AST does not attach as imp.Doc. By extending the byte range of the
	// adjacent import, the comment is preserved when LoadFormat() copies
	// src[d.Start:d.End] for each import.
	attachStandaloneComments(f.Comments, data, headEnd, tailStart, cStart, cEnd)

	sort.Sort(data)
	return data, headEnd, tailStart, cStart, cEnd, nil
}

// hasCommentBetweenDecls checks whether any comment group in the AST falls
// between two separate import declarations (not inside any of them).
func hasCommentBetweenDecls(comments []*ast.CommentGroup, decls []importDeclRange, cStart, cEnd int) bool {
	for _, cg := range comments {
		cgStart := int(cg.Pos()) - 1
		cgEnd := int(cg.End())

		// Skip comments inside the C import block.
		if cStart != 0 && cgStart >= cStart && cgEnd <= cEnd {
			continue
		}

		// Check if this comment falls between any two consecutive import declarations.
		// importDecls is in source order because f.Decls is iterated in source order.
		for i := 0; i < len(decls)-1; i++ {
			if cgStart >= decls[i].end && cgEnd <= decls[i+1].start {
				return true
			}
		}
	}
	return false
}

// attachStandaloneComments scans all comment groups from the AST and attaches
// standalone comments (those not already covered by any import's byte range)
// to the adjacent import by extending its Start or End.
//
// A standalone comment is one that falls within the import block boundaries
// (headEnd..tailStart) but is not covered by any existing GciImports entry.
// This typically happens when a comment is separated from the next import by
// a blank line, so the Go AST does not attach it as imp.Doc.
//
// The algorithm processes comments in source order. When a comment is attached
// to the next import by extending its Start, subsequent comments that fall
// within the now-extended range are automatically detected as "covered" and
// skipped. This correctly handles multiple consecutive standalone comments
// before a single import.
//
// Known limitation: if a standalone comment labels a section (e.g., "// Third-party imports")
// and is attached to the first import in that section, sorting within the section may move
// the comment to a non-first position. This is an inherent limitation of the "extend Start"
// approach — the comment travels with its attached import during sorting.
//
// Note: this function mutates data in place. The covered check for subsequent
// comments sees the updated Start values, which is essential for correctly
// handling multiple consecutive standalone comments before a single import.
func attachStandaloneComments(comments []*ast.CommentGroup, data ImportList, headEnd, tailStart, cStart, cEnd int) {
	if len(data) == 0 {
		return
	}

	// Sort imports by Start position for correct positional lookups.
	sort.Slice(data, func(i, j int) bool {
		return data[i].Start < data[j].Start
	})

	for _, cg := range comments {
		cgStart := int(cg.Pos()) - 1
		cgEnd := int(cg.End())

		// Skip comments outside the import block.
		if cgStart < headEnd || cgEnd > tailStart {
			continue
		}

		// Skip comments inside the C import block.
		if cStart != 0 && cgStart >= cStart && cgEnd <= cEnd {
			continue
		}

		// Check if this comment is already covered by an existing import's
		// byte range (e.g., attached as imp.Doc or imp.Comment by the AST).
		covered := false
		for _, imp := range data {
			if imp.Start <= cgStart && cgEnd <= imp.End {
				covered = true
				break
			}
		}
		if covered {
			continue
		}

		// This is a standalone comment. Find the next import after it and
		// extend that import's Start to include the comment bytes.
		attached := false
		for _, imp := range data {
			if imp.Start > cgStart {
				imp.Start = cgStart
				attached = true
				break
			}
		}

		if !attached {
			// Trailing comment after the last import — extend the last
			// import's End to include it.
			data[len(data)-1].End = cgEnd
		}
	}
}

// IsGeneratedFileByComment reports whether the source file is generated code.
// Using a bit laxer rules than https://golang.org/s/generatedcode to
// match more generated code.
// Taken from https://github.com/golangci/golangci-lint.
func IsGeneratedFileByComment(in string) bool {
	const (
		genCodeGenerated = "code generated"
		genDoNotEdit     = "do not edit"
		genAutoFile      = "autogenerated file"      // easyjson
		genAutoGenerated = "automatically generated" // genny
	)

	markers := []string{genCodeGenerated, genDoNotEdit, genAutoFile, genAutoGenerated}
	in = strings.ToLower(in)
	for _, marker := range markers {
		if strings.Contains(in, marker) {
			return true
		}
	}

	return false
}

type NoImportError struct{}

func (n NoImportError) Error() string {
	return "No imports"
}

func (i NoImportError) Is(err error) bool {
	_, ok := err.(NoImportError)
	return ok
}

// CommentBetweenImportsError is returned when standalone comments exist between
// separate import declarations. In this case, reformatting would merge the
// declarations and drop the comments, so we skip formatting entirely.
type CommentBetweenImportsError struct{}

func (c CommentBetweenImportsError) Error() string {
	return "standalone comment between separate import declarations"
}

func (c CommentBetweenImportsError) Is(err error) bool {
	_, ok := err.(CommentBetweenImportsError)
	return ok
}
