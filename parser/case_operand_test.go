package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// A prefix operator in the operand of the simple CASE form is legal in
// ClickHouse, which answers 7 for each statement below. The disambiguator
// read CASE itself as a column name, because an operator follows it.
func TestCaseOperandAcceptsPrefixOperator(t *testing.T) {
	cases := []string{
		"SELECT CASE -2.5 WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE -1 WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE -x WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE +1 WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE -toInt32(1) WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE NOT true WHEN 1 THEN 5 ELSE 7 END",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			statements, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
			require.Len(t, statements, 1)
		})
	}
}

// The other side of the same fork. A CASE with no WHEN is a column name.
// That reading is deliberate, thus the fix must keep it.
func TestCaseKeywordStaysAColumnName(t *testing.T) {
	cases := []string{
		"SELECT CASE",
		"SELECT CASE - 2.5",
		"SELECT CASE + 1",
		"SELECT CASE, other FROM t",
		"SELECT CASE AS c FROM t",
		"SELECT CASE FROM t",
		"SELECT a FROM t WHERE CASE > 0",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			statements, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
			require.Len(t, statements, 1)
		})
	}
}

// The forms that never regressed, so a later change cannot trade one of
// them for the fix.
func TestCaseOperandFormsStillParse(t *testing.T) {
	cases := []string{
		"SELECT CASE 2.5 WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE (-2.5) WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE x WHEN 1 THEN 5 ELSE 7 END",
		"SELECT CASE WHEN 1 THEN 5 ELSE 7 END",
		"SELECT quantile(0.5)((CASE 1 WHEN 1 THEN 5 ELSE 7 END))",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			statements, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
			require.Len(t, statements, 1)
		})
	}
}
