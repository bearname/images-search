package order

type PayStatus int

const (
	PayProcessing PayStatus = iota
	PaySuccess
	PayFailed
	PayCanceled
)
