package repl

import (
	"bufio"
	"cottagepie/evaluator"
	"cottagepie/lexer"
	"cottagepie/object"
	"cottagepie/parser"
	"fmt"
	"io"
)

const PROMPT = ">> "
const ERROR_MESSAGE = `
#========================#
#      ERRORS FOUND      #
#========================#
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	book := object.NewCookbook()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, book)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, ERROR_MESSAGE)
	io.WriteString(out, "\nEncountered errors:\n")

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
