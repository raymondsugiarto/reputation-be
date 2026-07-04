package migrate

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestStripSQLLineComments_RemovesLineComments(t *testing.T) {
	in := "-- top comment\nSELECT 1;\n-- trailing\n"
	got := stripSQLLineComments(in)
	if strings.Contains(got, "--") {
		t.Fatalf("expected no `--` in output, got %q", got)
	}
	if !strings.Contains(got, "SELECT 1") {
		t.Fatalf("expected statement body preserved, got %q", got)
	}
}

func TestStripSQLLineComments_HandlesMissingTrailingNewline(t *testing.T) {
	// A `--` comment at end-of-file with no trailing newline must not
	// swallow the rest of the script.
	in := "SELECT 1; -- tail comment"
	got := stripSQLLineComments(in)
	if !strings.Contains(got, "SELECT 1") {
		t.Fatalf("expected statement preserved, got %q", got)
	}
	if strings.Contains(got, "--") {
		t.Fatalf("expected comment stripped, got %q", got)
	}
}

func TestStripSQLLineComments_PreservesStringLiterals(t *testing.T) {
	// A `--` inside a string literal must NOT be treated as a comment.
	in := "INSERT INTO foo (note) VALUES ('a -- b');"
	got := stripSQLLineComments(in)
	if !strings.Contains(got, "'a -- b'") {
		t.Fatalf("expected string literal preserved, got %q", got)
	}
}

func TestSplitSQLStatements_BasicThreeStatements(t *testing.T) {
	in := "INSERT INTO a VALUES (1); INSERT INTO b VALUES (2); INSERT INTO c VALUES (3);"
	got := splitSQLStatements(in)
	if len(got) != 3 {
		t.Fatalf("expected 3 statements, got %d: %v", len(got), got)
	}
	for i, s := range got {
		if !strings.HasPrefix(strings.TrimSpace(s), "INSERT INTO") {
			t.Fatalf("statement %d wrong: %q", i, s)
		}
	}
}

func TestSplitSQLStatements_HonoursQuotes(t *testing.T) {
	in := "INSERT INTO t (note) VALUES ('hello; world'); SELECT 1;"
	got := splitSQLStatements(in)
	if len(got) != 2 {
		t.Fatalf("expected 2 statements, got %d: %v", len(got), got)
	}
	if !strings.Contains(got[0], "'hello; world'") {
		t.Fatalf("string literal must survive the split: %q", got[0])
	}
}

func TestIsDuplicateKeyError(t *testing.T) {
	cases := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{&fakeErr{msg: "ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)"}, true},
		{&fakeErr{msg: "Error 1062: Duplicate entry 'admin' for key 'username'"}, true},
		{&fakeErr{msg: "some other error"}, false},
	}
	for _, c := range cases {
		got := isDuplicateKeyError(c.err)
		if got != c.want {
			t.Errorf("isDuplicateKeyError(%v) = %v, want %v", c.err, got, c.want)
		}
	}
}

type fakeErr struct{ msg string }

func (e *fakeErr) Error() string { return e.msg }

// seedCustomerColumnsAllowList is the canonical column set for the
// `customer` table as produced by migrations
// 000001_init.up.sql + 000003_customer_type.up.sql +
// 000004_customer_company_fields.up.sql +
// 000005_customer_approval.up.sql.
//
// `shareholders` lives in 000004 (not 000005); `status`,
// `approved_by`, `approved_at`, `rejected_by`, `rejected_at`,
// `remark` were added in 000005.
//
// Every seed INSERT must reference columns from this set — anything
// outside it is a typo that will fail at runtime when the BE actually
// applies the seed. This test guards against the specific class of
// bug we just hit (`nama_pt` instead of `company_name`).
var seedCustomerColumnsAllowList = map[string]bool{
	// 000001
	"id":                   true,
	"organization_id":      true,
	"user_id":              true,
	"nama_lengkap":         true,
	"nomor_ktp":            true,
	"nomor_npwp":           true,
	"tanggal_lahir":        true,
	"kota_lahir":           true,
	"status_pernikahan":    true,
	"pendidikan_terakhir":  true,
	"lama_tinggal":         true,
	"alamat_jalan":         true,
	"kecamatan":            true,
	"kota_kabupaten":       true,
	"provinsi":             true,
	"kode_pos":             true,
	"sama_dengan_domisili": true,
	"alamat_ktp_jalan":     true,
	"kecamatan_ktp":        true,
	"kota_kabupaten_ktp":   true,
	"created_at":           true,
	"updated_at":           true,
	// 000003
	"customer_type": true,
	// 000004
	"company_name":       true,
	"establishment_date": true,
	"business_sector":    true,
	"kbli_code":          true,
	"company_tax_id":     true,
	"shareholders":       true,
	// 000005
	"status":      true,
	"approved_by": true,
	"approved_at": true,
	"rejected_by": true,
	"rejected_at": true,
	"remark":      true,
}

// insertIntoCustomerRegex captures the column list in
// `INSERT INTO customer (col1, col2, ...) VALUES (...)`.
// Multiline-aware so column lists that wrap across lines are matched
// in full.
var insertIntoCustomerRegex = regexp.MustCompile(`(?is)INSERT\s+INTO\s+customer\s*\(([^)]+)\)`)

// seedColumnLimits is the per-column VARCHAR ceiling for the
// customer table as declared by the migrations. `text` columns have
// no practical limit at the SQL level, so we model them with a very
// high sentinel (1 MB) — that's enough to flag obviously broken seed
// values like a stringified 5 MB blob.
//
// Keep this map in sync with the column types declared in
// db/migrations/postgres/000001_init.up.sql,
// 000004_customer_company_fields.up.sql, and
// 000005_customer_approval.up.sql. The shareholders column moved
// from varchar(255) to text in 000006_shareholders_text.up.sql.
var seedColumnLimits = map[string]int{
	"id":                   255,
	"organization_id":      255,
	"user_id":              255,
	"nama_lengkap":         255,
	"nomor_ktp":            255,
	"nomor_npwp":           255,
	"tanggal_lahir":        255, // TIMESTAMP — length check skipped
	"kota_lahir":           255,
	"status_pernikahan":    50,
	"pendidikan_terakhir":  50,
	"lama_tinggal":         255,
	"alamat_jalan":         255,
	"kecamatan":            255,
	"kota_kabupaten":       255,
	"provinsi":             255,
	"kode_pos":             20,
	"sama_dengan_domisili": 0, // boolean — skipped
	"alamat_ktp_jalan":     255,
	"kecamatan_ktp":        255,
	"kota_kabupaten_ktp":   255,
	"created_at":           255, // TIMESTAMP — length check skipped
	"updated_at":           255, // TIMESTAMP — length check skipped
	"customer_type":        50,
	"company_name":         255,
	"establishment_date":   255, // TIMESTAMP — length check skipped
	"business_sector":      255,
	"kbli_code":            255,
	"company_tax_id":       255,
	// After migration 000006 the shareholders column is `text`.
	// Sentinel of 1 MiB so we still catch obviously broken seed
	// values without flagging realistic JSON blobs.
	"shareholders":    1 << 20,
	"status":          50,
	"approved_by":     255,
	"approved_at":     255, // TIMESTAMP — length check skipped
	"rejected_by":     255,
	"rejected_at":     255, // TIMESTAMP — length check skipped
	"remark":          65535, // text column in 000005
}

func TestSeed_CustomerInsertsReferenceValidColumns(t *testing.T) {
	path := filepath.Join("..", "..", "..", "db", "migrations", "seed.sql")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read seed.sql: %v", err)
	}

	matches := insertIntoCustomerRegex.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		t.Fatal("expected at least one INSERT INTO customer in seed.sql")
	}

	for i, m := range matches {
		colsRaw := m[1]
		// The columns may span newlines; normalise to a comma list.
		cols := strings.Split(colsRaw, ",")
		for _, col := range cols {
			col = strings.TrimSpace(col)
			if col == "" {
				continue
			}
			if !seedCustomerColumnsAllowList[col] {
				t.Errorf("INSERT INTO customer #%d references unknown column %q — does it match the DB schema in db/migrations/postgres/00000*.up.sql?", i+1, col)
			}
		}
	}
}

// insertIntoCustomerFullRegex matches the whole INSERT statement so we
// can pair columns with their corresponding VALUES entries.
var insertIntoCustomerFullRegex = regexp.MustCompile(`(?is)INSERT\s+INTO\s+customer\s*\(([^)]+)\)\s*VALUES\s*\((.*?)\)\s*;`)

// unquoteSQL strips the surrounding single quotes from a SQL literal
// and decodes the doubled-quote escape ('' -> ').
func unquoteSQL(s string) string {
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return strings.ReplaceAll(s[1:len(s)-1], "''", "'")
	}
	return s
}

// splitTopLevelCommas splits s on commas that are not nested inside
// parens or single-quoted strings. Used to break a VALUES tuple into
// individual cell expressions.
func splitTopLevelCommas(s string) []string {
	var out []string
	var cur strings.Builder
	depth := 0
	inQ := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && (i == 0 || s[i-1] != '\\'):
			inQ = !inQ
			cur.WriteByte(c)
		case c == '(' && !inQ:
			depth++
			cur.WriteByte(c)
		case c == ')' && !inQ:
			depth--
			cur.WriteByte(c)
		case c == ',' && !inQ && depth == 0:
			out = append(out, cur.String())
			cur.Reset()
		default:
			cur.WriteByte(c)
		}
	}
	out = append(out, cur.String())
	for i := range out {
		out[i] = strings.TrimSpace(out[i])
	}
	return out
}

func TestSeed_CustomerInsertValuesRespectColumnLimits(t *testing.T) {
	path := filepath.Join("..", "..", "..", "db", "migrations", "seed.sql")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read seed.sql: %v", err)
	}

	matches := insertIntoCustomerFullRegex.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		t.Fatal("expected at least one INSERT INTO customer in seed.sql")
	}

	// Non-string columns: length check is meaningless, so we skip
	// these columns entirely.
	skipCheck := map[string]bool{
		"created_at":         true,
		"updated_at":         true,
		"tanggal_lahir":      true,
		"establishment_date": true,
		"approved_at":        true,
		"rejected_at":        true,
		"sama_dengan_domisili": true,
	}

	for i, m := range matches {
		cols := strings.Split(m[1], ",")
		for k := range cols {
			cols[k] = strings.TrimSpace(cols[k])
		}
		vals := splitTopLevelCommas(m[2])
		if len(vals) != len(cols) {
			t.Errorf("INSERT INTO customer #%d: column/value count mismatch (cols=%d, vals=%d)", i+1, len(cols), len(vals))
			continue
		}

		customerID := unquoteSQL(vals[0])
		for j, col := range cols {
			if skipCheck[col] {
				continue
			}
			limit, ok := seedColumnLimits[col]
			if !ok {
				// Column not in our limits table — falls through to
				// the column-allow-list test which would have failed
				// already.
				continue
			}

			raw := vals[j]
			// Strip ::cast suffix before measuring.
			if idx := strings.LastIndex(raw, "::"); idx >= 0 {
				raw = raw[:idx]
			}
			raw = strings.TrimSpace(raw)

			// Non-string literals (NOW(), true, false, numbers) are
			// skipped — their length is bounded by the literal text.
			if !strings.HasPrefix(raw, "'") {
				continue
			}
			inner := unquoteSQL(raw)
			// Count code points, not bytes — varchar(N) in Postgres
			// is character-count based, which matters for non-ASCII
			// like the Indonesian seed values.
			runeCount := len([]rune(inner))
			if runeCount > limit {
				t.Errorf(
					"INSERT INTO customer %s column %q value is %d chars (limit %d): %q",
					customerID, col, runeCount, limit,
					truncateForError(inner, 60),
				)
			}
		}
	}
}

func truncateForError(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

// Compile-time guard: assert that strconv is actually used so the
// import isn't dropped if we trim the test cases later.
var _ = strconv.Itoa
