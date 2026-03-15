package highlight

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var (
	formatter = html.New(
		html.WithClasses(true),
		html.PreventSurroundingPre(true),
	)
	// Use pygments style for maximum compatibility with existing Phorge CSS
	defaultStyle = styles.Get("pygments")
)

// Highlighter performs syntax highlighting using Chroma, producing
// Pygments-compatible CSS class names in <span> elements.
type Highlighter struct {
	lexerMap map[string]string
}

func New() *Highlighter {
	return &Highlighter{
		lexerMap: buildLexerMap(),
	}
}

type Result struct {
	HTML     string `json:"html"`
	Language string `json:"language"`
}

// Highlight takes source code and an optional language hint, returns
// HTML with <span class="..."> using Pygments-compatible CSS classes.
func (h *Highlighter) Highlight(source, language string) (*Result, error) {
	resolved := h.resolveLexer(language)

	lexer := lexers.Get(resolved)
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, defaultStyle, iterator); err != nil {
		return nil, err
	}

	output := buf.String()

	detectedLang := strings.ToLower(lexer.Config().Name)
	if language != "" {
		detectedLang = language
	}

	return &Result{
		HTML:     output,
		Language: detectedLang,
	}, nil
}

// Languages returns all supported lexer names.
func (h *Highlighter) Languages() []string {
	names := lexers.Names(false)
	result := make([]string, 0, len(names))
	for _, name := range names {
		result = append(result, strings.ToLower(name))
	}
	return result
}

func (h *Highlighter) resolveLexer(language string) string {
	if language == "" {
		return ""
	}
	lang := strings.ToLower(language)
	if mapped, ok := h.lexerMap[lang]; ok {
		return mapped
	}
	return lang
}

// buildLexerMap mirrors the Pygments lexer alias map from
// PhutilPygmentsSyntaxHighlighter::getPygmentsLexerNameFromLanguageName
// to ensure PHP sends the same language names and gets correct results.
func buildLexerMap() map[string]string {
	return map[string]string{
		"adb":             "ada",
		"ads":             "ada",
		"ahkl":            "ahk",
		"as":              "as3",
		"asax":            "aspx-vb",
		"ascx":            "aspx-vb",
		"ashx":            "aspx-vb",
		"asm":             "nasm",
		"asmx":            "aspx-vb",
		"aspx":            "aspx-vb",
		"autodelegate":    "myghty",
		"autohandler":     "mason",
		"aux":             "tex",
		"axd":             "aspx-vb",
		"b":               "brainfuck",
		"bas":             "vb.net",
		"bf":              "brainfuck",
		"bmx":             "blitzmax",
		"c++":             "cpp",
		"c++-objdump":     "cpp-objdump",
		"cc":              "cpp",
		"cfc":             "cfm",
		"cfg":             "ini",
		"cfml":            "cfm",
		"cl":              "common-lisp",
		"clj":             "clojure",
		"cmd":             "bat",
		"coffee":          "coffeescript",
		"cs":              "csharp",
		"csh":             "tcsh",
		"cw":              "redcode",
		"cxx":             "cpp",
		"cxx-objdump":     "cpp-objdump",
		"darcspatch":      "dpatch",
		"def":             "modula2",
		"dhandler":        "mason",
		"di":              "d",
		"duby":            "ruby",
		"dyl":             "dylan",
		"ebuild":          "bash",
		"eclass":          "bash",
		"el":              "common-lisp",
		"eps":             "postscript",
		"erl":             "erlang",
		"erl-sh":          "erl",
		"f":               "fortran",
		"f90":             "fortran",
		"feature":         "cucumber",
		"fhtml":           "velocity",
		"flx":             "felix",
		"flxh":            "felix",
		"frag":            "glsl",
		"g":               "antlr",
		"gdc":             "gooddata-cl",
		"gemspec":         "ruby",
		"geo":             "glsl",
		"gnumakefile":     "make",
		"h":               "c",
		"h++":             "cpp",
		"hh":              "cpp",
		"hpp":             "cpp",
		"hql":             "sql",
		"hrl":             "erlang",
		"hs":              "haskell",
		"htaccess":        "apacheconf",
		"htm":             "html",
		"hxx":             "cpp",
		"hy":              "hybris",
		"hyb":             "hybris",
		"ik":              "ioke",
		"inc":             "pov",
		"j":               "objective-j",
		"jbst":            "duel",
		"kid":             "genshi",
		"ksh":             "bash",
		"less":            "css",
		"lgt":             "logtalk",
		"lisp":            "common-lisp",
		"ll":              "llvm",
		"m":               "objective-c",
		"mak":             "make",
		"makefile":        "make",
		"man":             "groff",
		"mao":             "mako",
		"mc":              "mason",
		"md":              "minid",
		"mhtml":           "mason",
		"mi":              "mason",
		"ml":              "ocaml",
		"mli":             "ocaml",
		"mll":             "ocaml",
		"mly":             "ocaml",
		"mm":              "objective-c",
		"mo":              "modelica",
		"mod":             "modula2",
		"moo":             "moocode",
		"mu":              "mupad",
		"myt":             "myghty",
		"ns2":             "newspeak",
		"pas":             "delphi",
		"patch":           "diff",
		"phtml":           "html+php",
		"pl":              "prolog",
		"plot":            "gnuplot",
		"plt":             "gnuplot",
		"pm":              "perl",
		"po":              "pot",
		"pp":              "puppet",
		"pro":             "prolog",
		"proto":           "protobuf",
		"ps":              "postscript",
		"pxd":             "cython",
		"pxi":             "cython",
		"py":              "python",
		"pyw":             "python",
		"pyx":             "cython",
		"r":               "rebol",
		"r3":              "rebol",
		"rake":            "ruby",
		"rakefile":        "ruby",
		"rbw":             "ruby",
		"rbx":             "ruby",
		"rest":            "rst",
		"rl":              "ragel",
		"robot":           "robotframework",
		"rout":            "rconsole",
		"rss":             "xml",
		"s":               "gas",
		"sc":              "python",
		"scm":             "scheme",
		"sconscript":      "python",
		"sconstruct":      "python",
		"scss":            "css",
		"sh":              "bash",
		"sh-session":      "console",
		"spt":             "cheetah",
		"sqlite3-console": "sqlite3",
		"st":              "smalltalk",
		"sv":              "verilog",
		"tac":             "python",
		"tmpl":            "cheetah",
		"toc":             "tex",
		"tpl":             "smarty",
		"txt":             "text",
		"vapi":            "vala",
		"vb":              "vb.net",
		"vert":            "glsl",
		"vhd":             "vhdl",
		"vimrc":           "vim",
		"vm":              "velocity",
		"weechatlog":      "irc",
		"wlua":            "lua",
		"wsdl":            "xml",
		"xhtml":           "html",
		"xqy":             "xquery",
		"xsd":             "xml",
		"xsl":             "xslt",
		"xslt":            "xml",
		"yml":             "yaml",
		"rs":              "rust",
		"ts":              "typescript",
		"tsx":             "typescriptreact",
		"jsx":             "jsx",
		"kt":              "kotlin",
		"kts":             "kotlin",
		"swift":           "swift",
		"gradle":          "groovy",
		"tf":              "terraform",
		"hcl":             "terraform",
		"dockerfile":      "docker",
		"containerfile":   "docker",
		"toml":            "toml",
		"graphql":         "graphql",
		"gql":             "graphql",
	}
}
