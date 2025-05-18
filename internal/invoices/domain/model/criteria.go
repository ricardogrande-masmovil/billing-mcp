package model

import "time"

type Criteria struct {
	Status        InvoiceStatus
	IssueDateFrom time.Time
	IssueDateTo   time.Time
}
