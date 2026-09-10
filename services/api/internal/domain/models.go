package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeCash       AccountType = "cash"
	AccountTypeInvestment AccountType = "investment"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleViewer Role = "viewer"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Household struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	Locale    string    `json:"locale"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type HouseholdMember struct {
	ID          uuid.UUID `json:"id"`
	HouseholdID uuid.UUID `json:"householdId"`
	UserID      uuid.UUID `json:"userId"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
}

type FinancialAccount struct {
	ID                  uuid.UUID   `json:"id"`
	HouseholdID         uuid.UUID   `json:"householdId"`
	Name                string      `json:"name"`
	Type                AccountType `json:"type"`
	Currency            string      `json:"currency"`
	InitialBalanceMinor int64       `json:"initialBalanceMinor"`
	CurrentBalanceMinor int64       `json:"currentBalanceMinor"`
	CreditLimitMinor    int64       `json:"creditLimitMinor,omitempty"`
	ClosingDay          int         `json:"statementClosingDay,omitempty"`
	DueDay              int         `json:"statementDueDay,omitempty"`
	IsActive            bool        `json:"isActive"`
	CreatedAt           time.Time   `json:"createdAt"`
	UpdatedAt           time.Time   `json:"updatedAt"`
}

func (a FinancialAccount) FormattedCurrentBalance() string {
	return NewMoney(a.CurrentBalanceMinor, a.Currency).FormatBRL()
}

type Category struct {
	ID          uuid.UUID  `json:"id"`
	HouseholdID uuid.UUID  `json:"householdId"`
	Name        string     `json:"name"`
	Icon        string     `json:"icon"`
	Color       string     `json:"color"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type AuditEvent struct {
	ID           uuid.UUID              `json:"id"`
	HouseholdID  uuid.UUID              `json:"householdId"`
	ActorID      *uuid.UUID             `json:"actorId,omitempty"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resourceType"`
	ResourceID   *uuid.UUID             `json:"resourceId,omitempty"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
}
