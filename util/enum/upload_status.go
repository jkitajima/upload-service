package enum

type UploadStatus uint

const (
	NotAvailable UploadStatus = iota
	Pending
	Processing
	Retrying
	Failed
	Canceled
	Completed
)

func (s UploadStatus) String() string {
	switch s {
	case NotAvailable:
		return "Not available"
	case Pending:
		return "Pending"
	case Processing:
		return "Processing"
	case Retrying:
		return "Retrying"
	case Failed:
		return "Failed"
	case Canceled:
		return "Canceled"
	case Completed:
		return "Completed"
	default:
		return "Not available"
	}
}
