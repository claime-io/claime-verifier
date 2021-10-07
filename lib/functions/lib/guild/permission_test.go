package guild

import (
	"testing"
)

func TestPermission(t *testing.T) {
	t.Run("has Permission", func(t *testing.T) {
		if !HasPermissionAdministrator(1099511627775) {
			t.Error("should be true")
		}
	})
	t.Run("does not have Permission", func(t *testing.T) {
		if HasPermissionAdministrator(1071698660928) {
			t.Error("should be false")
		}
	})
}
