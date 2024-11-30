package entity

type DisbursementStatusEnum uint8

const (
	DisbursementStatusUnspecified DisbursementStatusEnum = iota
	DisbursementStatusPending
	DisbursementStatusSuccess
	DisbursementStatusFailed
)

func (status DisbursementStatusEnum) String() string {
	switch status {
	case DisbursementStatusPending:
		return "PENDING"
	case DisbursementStatusFailed:
		return "FAILED"
	case DisbursementStatusSuccess:
		return "SUCCESS"
	default:
		return "UNSPECIFIED"
	}
}
