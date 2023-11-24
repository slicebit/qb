package postgres

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/orus-io/qb"
)

// Dialect is a type of dialect that can be used with postgres driver
type Dialect struct {
	bindingIndex int
	escaping     bool
}

// NewDialect returns a new PostgresDialect
func NewDialect() qb.Dialect {
	return &Dialect{escaping: false, bindingIndex: 0}
}

func init() {
	qb.RegisterDialect("postgres", NewDialect())
}

// CompileType compiles a type into its DDL
func (d *Dialect) CompileType(t qb.TypeElem) string {
	if t.Name == "BLOB" {
		return "bytea"
	}
	return qb.DefaultCompileType(t, d.SupportsUnsigned())
}

// Escape wraps the string with escape characters of the dialect
func (d *Dialect) Escape(str string) string {
	if d.escaping {
		return fmt.Sprintf("\"%s\"", str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (d *Dialect) EscapeAll(strings []string) []string {
	return qb.EscapeAll(d, strings[0:])
}

// SetEscaping sets the escaping parameter of dialect
func (d *Dialect) SetEscaping(escaping bool) {
	d.escaping = escaping
}

// Escaping gets the escaping parameter of dialect
func (d *Dialect) Escaping() bool {
	return d.escaping
}

// AutoIncrement generates auto increment sql of current dialect
func (d *Dialect) AutoIncrement(column *qb.ColumnElem) string {
	var colSpec string
	if column.Type.Name == "BIGINT" {
		colSpec = "BIGSERIAL"
	} else if column.Type.Name == "SMALLINT" {
		colSpec = "SMALLSERIAL"
	} else {
		colSpec = "SERIAL"
	}
	if column.Options.InlinePrimaryKey {
		colSpec += " PRIMARY KEY"
	}
	return colSpec
}

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *Dialect) SupportsUnsigned() bool { return false }

// Driver returns the current driver of dialect
func (d *Dialect) Driver() string {
	return "postgres"
}

// GetCompiler returns a PostgresCompiler
func (d *Dialect) GetCompiler() qb.Compiler {
	return PostgresCompiler{qb.NewSQLCompiler(d)}
}

// WrapError wraps a native error in a qb Error
func (d *Dialect) WrapError(err error) (qbErr qb.Error) {
	qbErr.Orig = err
	pgErr, ok := err.(*pq.Error)
	if !ok {
		return
	}
	switch pgErr.Code.Class() {
	case "0A": // Class 0A - Feature Not Supported
		qbErr.Code = qb.ErrNotSupported
	case "20", // Class 20 - Case Not Found
		"21": //  Class 21 - Cardinality Violation
		qbErr.Code = qb.ErrProgramming
	case "22": // Class 22 - Data Exception
		qbErr.Code = qb.ErrData
	case "23": // Class 23 - Integrity Constraint Violation
		qbErr.Code = qb.ErrIntegrity
	case "24", // Class 24 - Invalid Cursor State
		"25": //  Class 25 - Invalid Transaction State
		qbErr.Code = qb.ErrInternal
	case "26", // Class 26 - Invalid SQL Statement Name
		"27", //  Class 27 - Triggered Data Change Violation
		"28": //  Class 28 - Invalid Authorization Specification
		qbErr.Code = qb.ErrOperational
	case "2B", // Class 2B - Dependent Privilege Descriptors Still Exist
		"2D", //  Class 2D - Invalid Transaction Termination
		"2F": //  Class 2F - SQL Routine Exception
		qbErr.Code = qb.ErrInternal
	case "34": // Class 34 - Invalid Cursor Name
		qbErr.Code = qb.ErrOperational
	case "38", // Class 38 - External Routine Exception
		"39", //  Class 39 - External Routine Invocation Exception
		"3B": //  Class 3B - Savepoint Exception
		qbErr.Code = qb.ErrInternal
	case "3D", // Class 3D - Invalid Catalog Name
		"3F": //  Class 3F - Invalid Schema Name
		qbErr.Code = qb.ErrProgramming
	case "40": // Class 40 - Transaction Rollback
		qbErr.Code = qb.ErrOperational
	case "42", // Class 42 - Syntax Error or Access Rule Violation
		"44": //  Class 44 - WITH CHECK OPTION Violation
		qbErr.Code = qb.ErrProgramming
	case "53", // Class 53 - Insufficient Resources
		"54", //  Class 54 - Program Limit Exceeded
		"55", //  Class 55 - Object Not In Prerequisite State
		"57", //  Class 57 - Operator Intervention
		"58": //  Class 58 - System Error (errors external to PostgreSQL itself)
		qbErr.Code = qb.ErrOperational

	case "F0": // Class F0 - Configuration File Error
		qbErr.Code = qb.ErrInternal
	case "HV": // Class HV - Foreign Data Wrapper Error (SQL/MED)
		qbErr.Code = qb.ErrOperational
	case "P0", // Class P0 - PL/pgSQL Error
		"XX": //  Class XX - Internal Error
		qbErr.Code = qb.ErrInternal
	default:
		qbErr.Code = qb.ErrDatabase
	}
	return
}

// PostgresCompiler is a SQLCompiler specialised for PostgreSQL
type PostgresCompiler struct {
	qb.SQLCompiler
}

// VisitBind renders a bounded value
func (PostgresCompiler) VisitBind(context *qb.CompilerContext, bind qb.BindClause) string {
	context.Binds = append(context.Binds, bind.Value)
	return fmt.Sprintf("$%d", len(context.Binds))
}

// VisitUpsert generates INSERT INTO ... VALUES ... ON CONFLICT(...) DO UPDATE SET ...
func (PostgresCompiler) VisitUpsert(context *qb.CompilerContext, upsert qb.UpsertStmt) string {
	var (
		colNames []string
		values   []string
	)
	for k, v := range upsert.ValuesMap {
		colNames = append(colNames, context.Compiler.VisitLabel(context, k))
		context.Binds = append(context.Binds, v)
		values = append(values, fmt.Sprintf("$%d", len(context.Binds)))
	}

	var updates []string
	for k, v := range upsert.ValuesMap {
		context.Binds = append(context.Binds, v)
		updates = append(updates, fmt.Sprintf(
			"%s = %s",
			context.Dialect.Escape(k),
			fmt.Sprintf("$%d", len(context.Binds)),
		))
	}

	var uniqueCols []string
	for _, c := range upsert.Table.PrimaryCols() {
		uniqueCols = append(uniqueCols, context.Compiler.VisitLabel(context, c.Name))
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s(%s)\nVALUES(%s)\nON CONFLICT (%s) DO UPDATE SET %s",
		context.Compiler.VisitLabel(context, upsert.Table.Name),
		strings.Join(colNames, ", "),
		strings.Join(values, ", "),
		strings.Join(uniqueCols, ", "),
		strings.Join(updates, ", "))

	var returning []string
	for _, r := range upsert.ReturningCols {
		returning = append(returning, context.Compiler.VisitLabel(context, r.Name))
	}
	if len(returning) > 0 {
		sql += fmt.Sprintf(
			"RETURNING %s",
			strings.Join(returning, ", "),
		)
	}
	return sql
}
