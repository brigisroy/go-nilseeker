package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// checkSelectorExpr checks for nil pointer dereferences in selector expressions (e.g., a.b)
func checkSelectorExpr(pass *analysis.Pass, node *ast.SelectorExpr, checkedVars map[string]bool) {
	// Get the type information for the expression
	exprType := pass.TypesInfo.Types[node.X]
	if exprType.Type == nil {
		return // No type information available
	}

	// If it's a pointer type, check if it might be nil
	if _, ok := exprType.Type.Underlying().(*types.Pointer); ok {
		// Check if it's a variable and if it's been checked for nil
		if ident, ok := node.X.(*ast.Ident); ok {
			if !checkedVars[ident.Name] {
				// Report potential nil dereference
				pass.Reportf(node.Pos(), "potential nil pointer dereference: %s may be nil", ident.Name)
			}
		} else {
			// It's not a simple variable, so it might be something like a function return
			// that hasn't been checked
			pass.Reportf(node.Pos(), "potential nil pointer dereference in selector expression")
		}
	}
}

// checkStarExpr checks explicit pointer dereferences (e.g., *p)
func checkStarExpr(pass *analysis.Pass, node *ast.StarExpr, checkedVars map[string]bool) {
	// Check if the expression being dereferenced might be nil
	if ident, ok := node.X.(*ast.Ident); ok {
		if !checkedVars[ident.Name] {
			pass.Reportf(node.Pos(), "explicit dereference of possibly nil pointer: %s", ident.Name)
		}
	} else {
		// It's a more complex expression
		pass.Reportf(node.Pos(), "explicit dereference of possibly nil pointer")
	}
}

// checkIndexExpr checks for potential nil slice or map indexing
func checkIndexExpr(pass *analysis.Pass, node *ast.IndexExpr, checkedVars map[string]bool) {
	// Get type information for the expression being indexed
	exprType := pass.TypesInfo.Types[node.X]
	if exprType.Type == nil {
		return // No type information available
	}

	underlying := exprType.Type.Underlying()

	// Check if it's a slice or map
	switch underlying.(type) {
	case *types.Slice, *types.Map:
		if ident, ok := node.X.(*ast.Ident); ok {
			if !checkedVars[ident.Name] {
				pass.Reportf(node.Pos(), "indexing potentially nil %s: %s", typeKind(underlying), ident.Name)
			}
		} else {
			pass.Reportf(node.Pos(), "indexing potentially nil %s", typeKind(underlying))
		}
	}
}

// typeKind returns a string description of the type for error messages
func typeKind(t types.Type) string {
	switch t.(type) {
	case *types.Slice:
		return "slice"
	case *types.Map:
		return "map"
	case *types.Pointer:
		return "pointer"
	default:
		return "value"
	}
}

// trackNilChecks looks for nil checks in if statements and marks variables as checked
func trackNilChecks(node *ast.IfStmt, checkedVars map[string]bool) {
	// Pattern: if x != nil { ... }
	if binExpr, ok := node.Cond.(*ast.BinaryExpr); ok {
		if binExpr.Op == token.NEQ {
			// Check if one side is nil and the other is an identifier
			var ident *ast.Ident

			// Left side is nil, right side is variable
			if isNilIdent(binExpr.X) && isIdent(binExpr.Y) {
				ident = binExpr.Y.(*ast.Ident)
			}

			// Right side is nil, left side is variable
			if isNilIdent(binExpr.Y) && isIdent(binExpr.X) {
				ident = binExpr.X.(*ast.Ident)
			}

			// Mark the variable as checked
			if ident != nil {
				checkedVars[ident.Name] = true
			}
		}
	}
}

// isNilIdent checks if an expression is the nil identifier
func isNilIdent(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "nil"
	}
	return false
}

// isIdent checks if an expression is an identifier
func isIdent(expr ast.Expr) bool {
	_, ok := expr.(*ast.Ident)
	return ok
}
