package entity_test

import (
	"expense-management-system/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRole_ParseUserRole(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		wantRes    entity.UserRole
		wantErrMsg string
	}{
		{
			name:       "employee role",
			status:     "employee",
			wantRes:    entity.UserRoleEmployee,
			wantErrMsg: "",
		},
		{
			name:       "manager role",
			status:     "manager",
			wantRes:    entity.UserRoleManager,
			wantErrMsg: "",
		},
		{
			name:       "unknown role",
			status:     "unknown",
			wantRes:    "",
			wantErrMsg: "invalid user role = unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := entity.ParseUserRole(tt.status)

			assert.Equal(t, tt.wantRes, res)
			if tt.wantErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.wantErrMsg, err.Error())
			}
		})
	}
}
