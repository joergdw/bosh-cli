package protogetter

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

type Mode int

const (
	StandaloneMode Mode = iota
	GolangciLintMode
)

const msgFormat = "avoid direct access to proto field %s, use %s instead"

func NewAnalyzer(cfg *Config) *analysis.Analyzer {
	if cfg == nil {
		cfg = &Config{}
	}

	return &analysis.Analyzer{
		Name:  "protogetter",
		Doc:   "Reports direct reads from proto message fields when getters should be used",
		Flags: flags(cfg),
		Run: func(pass *analysis.Pass) (any, error) {
			_, err := Run(pass, cfg)
			return nil, err
		},
	}
}

func flags(opts *Config) flag.FlagSet {
	fs := flag.NewFlagSet("protogetter", flag.ContinueOnError)

	fs.Func("skip-generated-by", "skip files generated with the given prefixes", func(s string) error {
		for _, prefix := range strings.Split(s, ",") {
			opts.SkipGeneratedBy = append(opts.SkipGeneratedBy, prefix)
		}
		return nil
	})
	fs.Func("skip-files", "skip files with the given glob patterns", func(s string) error {
		for _, pattern := range strings.Split(s, ",") {
			opts.SkipFiles = append(opts.SkipFiles, pattern)
		}
		return nil
	})
	fs.BoolVar(&opts.SkipAnyGenerated, "skip-any-generated", false, "skip any generated files")

	return *fs
}

type Config struct {
	Mode             Mode // Zero value is StandaloneMode.
	SkipGeneratedBy  []string
	SkipFiles        []string
	SkipAnyGenerated bool
}

func Run(pass *analysis.Pass, cfg *Config) ([]Issue, error) {
	skipGeneratedBy := make([]string, 0, len(cfg.SkipGeneratedBy)+3)
	// Always skip files generated by protoc-gen-go, protoc-gen-go-grpc and protoc-gen-grpc-gateway.
	skipGeneratedBy = append(skipGeneratedBy, "protoc-gen-go", "protoc-gen-go-grpc", "protoc-gen-grpc-gateway")
	for _, s := range cfg.SkipGeneratedBy {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		skipGeneratedBy = append(skipGeneratedBy, s)
	}

	skipFilesGlobPatterns := make([]glob.Glob, 0, len(cfg.SkipFiles))
	for _, s := range cfg.SkipFiles {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		compile, err := glob.Compile(s)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern: %w", err)
		}

		skipFilesGlobPatterns = append(skipFilesGlobPatterns, compile)
	}

	nodeTypes := []ast.Node{
		(*ast.AssignStmt)(nil),
		(*ast.CallExpr)(nil),
		(*ast.SelectorExpr)(nil),
		(*ast.IncDecStmt)(nil),
		(*ast.UnaryExpr)(nil),
	}

	// Skip filtered files.
	var files []*ast.File
	for _, f := range pass.Files {
		if skipGeneratedFile(f, skipGeneratedBy, cfg.SkipAnyGenerated) {
			continue
		}

		if skipFilesByGlob(pass.Fset.File(f.Pos()).Name(), skipFilesGlobPatterns) {
			continue
		}

		files = append(files, f)

		// ast.Print(pass.Fset, f)
	}

	ins := inspector.New(files)

	var issues []Issue

	filter := NewPosFilter()
	ins.Preorder(nodeTypes, func(node ast.Node) {
		report := analyse(pass, filter, node)
		if report == nil {
			return
		}

		switch cfg.Mode {
		case StandaloneMode:
			pass.Report(report.ToDiagReport())
		case GolangciLintMode:
			issues = append(issues, report.ToIssue(pass.Fset))
		}
	})

	return issues, nil
}

func analyse(pass *analysis.Pass, filter *PosFilter, n ast.Node) *Report {
	// fmt.Printf("\n>>> check: %s\n", formatNode(n))
	// ast.Print(pass.Fset, n)
	if filter.IsFiltered(n.Pos()) {
		// fmt.Printf(">>> filtered\n")
		return nil
	}

	result, err := Process(pass.TypesInfo, filter, n)
	if err != nil {
		pass.Report(analysis.Diagnostic{
			Pos:     n.Pos(),
			End:     n.End(),
			Message: fmt.Sprintf("error: %v", err),
		})

		return nil
	}

	// If existing in filter, skip it.
	if filter.IsFiltered(n.Pos()) {
		return nil
	}

	if result.Skipped() {
		return nil
	}

	// If the expression has already been replaced, skip it.
	if filter.IsAlreadyReplaced(pass.Fset, n.Pos(), n.End()) {
		return nil
	}
	// Add the expression to the filter.
	filter.AddAlreadyReplaced(pass.Fset, n.Pos(), n.End())

	return &Report{
		node:   n,
		result: result,
	}
}

// Issue is used to integrate with golangci-lint's inline auto fix.
type Issue struct {
	Pos       token.Position
	Message   string
	InlineFix InlineFix
}

type InlineFix struct {
	StartCol  int // zero-based
	Length    int
	NewString string
}

type Report struct {
	node   ast.Node
	result *Result
}

func (r *Report) ToDiagReport() analysis.Diagnostic {
	msg := fmt.Sprintf(msgFormat, r.result.From, r.result.To)

	return analysis.Diagnostic{
		Pos:     r.node.Pos(),
		End:     r.node.End(),
		Message: msg,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: msg,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     r.node.Pos(),
						End:     r.node.End(),
						NewText: []byte(r.result.To),
					},
				},
			},
		},
	}
}

func (r *Report) ToIssue(fset *token.FileSet) Issue {
	msg := fmt.Sprintf(msgFormat, r.result.From, r.result.To)
	return Issue{
		Pos:     fset.Position(r.node.Pos()),
		Message: msg,
		InlineFix: InlineFix{
			StartCol:  fset.Position(r.node.Pos()).Column - 1,
			Length:    len(r.result.From),
			NewString: r.result.To,
		},
	}
}

func skipGeneratedFile(f *ast.File, prefixes []string, skipAny bool) bool {
	if len(f.Comments) == 0 {
		return false
	}
	firstComment := f.Comments[0].Text()

	// https://golang.org/s/generatedcode
	if skipAny && strings.HasPrefix(firstComment, "Code generated") {
		return true
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(firstComment, "Code generated by "+prefix) {
			return true
		}
	}

	return false
}

func skipFilesByGlob(filename string, patterns []glob.Glob) bool {
	for _, pattern := range patterns {
		if pattern.Match(filename) || pattern.Match(filepath.Base(filename)) {
			return true
		}
	}

	return false
}

func formatNode(node ast.Node) string {
	buf := new(bytes.Buffer)
	if err := format.Node(buf, token.NewFileSet(), node); err != nil {
		log.Printf("Error formatting expression: %v", err)
		return ""
	}

	return buf.String()
}
