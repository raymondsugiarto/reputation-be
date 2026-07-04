package migrate

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/raymondsugiarto/reputation-be/config"
	"gorm.io/gorm"
)

// seedFilePath is the canonical seed file shipped alongside the
// migrations. The path is relative to the working directory because
// the rest of the migrate package also resolves paths that way —
// keeping them symmetric avoids the "works in dev, broken in Docker"
// trap.
const seedFilePath = "db/migrations/seed.sql"

// Seed applies the project seed file against the configured database.
// It is safe to call repeatedly: every INSERT uses fixed primary keys
// (e.g. user id '1', customer id 'seed-cust-pending-1'), so the second
// run produces a primary-key violation that we detect and treat as a
// no-op ("already seeded"). This lets `MigrateUpAll` call Seed on
// every production boot without manual state tracking.
//
// Returns the number of statements successfully executed. Errors that
// indicate the file is already loaded (PK conflicts, duplicate seed
// rows) are converted into a warning so the boot can continue.
func Seed(db *gorm.DB) error {
	if db == nil {
		return errors.New("seed: nil db")
	}

	cfg := config.GetConfig().Database.Main
	schema := strings.TrimSpace(cfg.Schema)
	if schema != "" && schema != "public" {
		// The seed file doesn't parameterise schema — it inserts into
		// the public schema. We honour the configured search_path so
		// the inserts land in the right place.
		if err := db.Exec(fmt.Sprintf("SET search_path TO %q", schema)).Error; err != nil {
			return fmt.Errorf("seed: set search_path: %w", err)
		}
	}

	body, err := os.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf("seed: read %s: %w", seedFilePath, err)
	}

	// Strip line comments. The seed file is hand-written SQL — it
	// uses `--` line comments which GORM's Exec does not strip, and
	// some Postgres drivers choke on multi-line `--` blocks mid-stmt.
	bodyStr := stripSQLLineComments(string(body))

	// Split on `;` boundaries. The naive split is fine for this seed
	// file (no semicolons inside string literals or dollar-quoted
	// blocks); the production migrations folder is the same shape.
	statements := splitSQLStatements(bodyStr)

	executed := 0
	skipped := 0
	for _, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		if err := db.Exec(trimmed).Error; err != nil {
			// Already-seeded errors. Postgres returns SQLSTATE 23505
			// (unique_violation) on duplicate primary keys. We treat
			// that as "skip this statement" so re-running Seed is a
			// no-op rather than a fatal error.
			if isDuplicateKeyError(err) {
				skipped++
				continue
			}
			return fmt.Errorf("seed: stmt %d: %w", executed+1, err)
		}
		executed++
	}
	fmt.Printf("seed: applied %d statements (%d skipped as already-seeded)\n", executed, skipped)
	return nil
}

// stripSQLLineComments removes `-- …\n` comments. We deliberately keep
// block comments (`/* … */`) untouched because the seed file doesn't
// use them and stripping them safely requires a real SQL parser.
//
// `--` inside a single-quoted string literal is preserved. The same
// literal must already be protected from the statement splitter
// (see splitSQLStatements), so we don't need to honour quotes here —
// but we do, as a defensive measure.
func stripSQLLineComments(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			inQuote = !inQuote
			b.WriteByte(c)
			continue
		}
		if !inQuote && c == '-' && i+1 < len(s) && s[i+1] == '-' {
			// Walk forward to the next newline (or `;` as a fallback
			// so we don't consume the rest of the file if the comment
			// is the last line of the file and has no trailing \n).
			j := i
			for j < len(s) && s[j] != '\n' && s[j] != ';' {
				j++
			}
			// Re-emit the newline / semicolon if we hit one so
			// downstream splitting still works correctly.
			if j < len(s) {
				b.WriteByte(s[j])
				i = j
			} else {
				i = j - 1
			}
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

// splitSQLStatements splits a SQL script on top-level semicolons.
// Single-quoted strings are honoured so semicolons inside literals
// don't split the script.
func splitSQLStatements(s string) []string {
	var (
		out     []string
		current strings.Builder
		inQuote bool
	)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'':
			inQuote = !inQuote
			current.WriteByte(c)
		case c == ';' && !inQuote:
			out = append(out, current.String())
			current.Reset()
		default:
			current.WriteByte(c)
		}
	}
	if strings.TrimSpace(current.String()) != "" {
		out = append(out, current.String())
	}
	return out
}

// isDuplicateKeyError matches Postgres unique_violation (23505) and the
// MySQL equivalent (1062). The seed file uses both schemas because
// the platform can run on either — the migration runner picks one at
// build time via the driver import. Detecting both lets us run the
// same seed regardless of the chosen adapter.
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.Contains(msg, "23505") || strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "1062") {
		return true
	}
	return false
}
