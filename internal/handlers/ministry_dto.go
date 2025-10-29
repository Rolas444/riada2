package handlers

type CreateMinistryRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Mission     *string `json:"mission"`
	Status      string  `json:"status"`
}

type UpdateMinistryRequest struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Mission     *string `json:"mission"`
	Status      string  `json:"status"`
}

type CreateMinistryMemberRequest struct {
	MinistryID uint    `json:"ministryID"`
	PersonID   uint    `json:"personID"`
	Role       *string `json:"role"`
	Status     string  `json:"status"`
}

type UpdateMinistryMemberRequest struct {
	MinistryID uint    `json:"ministryID"`
	PersonID   uint    `json:"personID"`
	Role       *string `json:"role"`
	Status     string  `json:"status"`
}
