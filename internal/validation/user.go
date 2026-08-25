package validation

import (
	"fmt"
	"strings"
	"timber-safety/internal/domain"
)

func ValidateUser(u domain.User) error {
	if strings.TrimSpace(u.ID) == "" || strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("user id and name are required")
	}
	if u.Role != "registrar" && u.Role != "inspector" && u.Role != "manager" {
		return fmt.Errorf("unsupported role")
	}
	return nil
}
func CanReview(u domain.User) bool { return u.Active && (u.Role == "inspector" || u.Role == "manager") }
func CanRegister(u domain.User) bool {
	return u.Active && (u.Role == "registrar" || u.Role == "manager")
}
func NormalizeRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	switch role {
	case "登记员":
		return "registrar"
	case "审核员":
		return "inspector"
	case "管理员":
		return "manager"
	}
	return role
}
