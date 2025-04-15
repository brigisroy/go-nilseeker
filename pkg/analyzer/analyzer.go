// Package analyzer provides functionality for detecting potential nil pointer dereferences in Go code.
package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NilPointerAnalyzer returns an analyzer that checks for likely nil pointer dereferences in Go code.
func NilPointerAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "nilseeker",
		Doc:      "Detects potential nil pointer dereferences",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
}

// run implements the analysis logic for nil pointer detection
func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Define which node types we're interested in
	nodeFilter := []ast.Node{
		(*ast.SelectorExpr)(nil), // For field access/method calls (obj.field, obj.Method())
		(*ast.StarExpr)(nil),     // For explicit dereferences (*ptr)
		(*ast.IndexExpr)(nil),    // For slice/map indexing (slice[i], map[key])
		(*ast.CallExpr)(nil),     // For function calls that might return nil
		(*ast.AssignStmt)(nil),   // For assignments that might assign nil
		(*ast.IfStmt)(nil),       // For tracking nil checks
	}

	// Track variables that have been checked for nil
	// This simple version doesn't account for scope or control flow
	checkedVars := make(map[string]bool)

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			checkSelectorExpr(pass, node, checkedVars)
		case *ast.StarExpr:
			checkStarExpr(pass, node, checkedVars)
		case *ast.IndexExpr:
			checkIndexExpr(pass, node, checkedVars)
		case *ast.IfStmt:
			trackNilChecks(node, checkedVars)
		}
	})

	return nil, nil
}
