package entity

type DisbursementStatus uint8

const (
	DisbursementStatusUnspecified DisbursementStatus = iota
	DisbursementStatusPending
	DisbursementStatusSuccess
	DisbursementStatusFailed
)
