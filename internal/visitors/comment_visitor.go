package visitors

import (
	"github.com/pm1381/sirish/internal/dto"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Comment struct {
	targetInterfaces dto.Types
	fileAbsPath      string // absolutePath
}

func NewCommentVisitor(filePath string) *Comment {
	return &Comment{
		fileAbsPath: filePath,
	}
}

func (c *Comment) Traverse() {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, c.fileAbsPath, nil, parser.ParseComments)
	// you can use file name and file contents too as src.
	if err != nil {
		panic(err)
	}
	ast.Walk(c, file)
}

func (c *Comment) GetTargets() dto.Types {
	return c.targetInterfaces
}

func (c *Comment) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil // reached the end in depth traverse
	}
	switch nodeWithType := node.(type) {
	case *ast.Comment:
		if strings.Contains(nodeWithType.Text, "// sirish:") {
			c.targetInterfaces = append(c.targetInterfaces, strings.TrimSpace(strings.Split(nodeWithType.Text, "// sirish:")[1]))
		} else if strings.Contains(nodeWithType.Text, "//sirish:") {
			c.targetInterfaces = append(c.targetInterfaces, strings.TrimSpace(strings.Split(nodeWithType.Text, "//sirish:")[1]))
		}
	}
	return c
}
