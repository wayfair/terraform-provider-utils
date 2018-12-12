// Package validation contains schema.Schema validation functions. This package
// is designed similarly to the terraform/helper/validation package.
package validation

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// DiffSuppressStringIgnoreCase suppresses a Resource's schema diff for a
// Schema.TypeString attribute if the old and new version of the string are the
// same when ignoring case.
//
// k contains the schema's "key". This is a value in the format of a terraform
// tfstate file.  old contains the previously known value for that key. new
// contains the new value for that key.
//
// Returns true if the diff should be suppressed, false if the diff should
// be evaluated.
func DiffSuppressStringIgnoreCase(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
