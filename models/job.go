package models

type Job struct {
	Title    string
	Location string
	Link     string
	Status   JobStatus
}

type JobStatus string

const (
	TO_APPLY     JobStatus = "🔵 TO_APPLY"
	APPLIED      JobStatus = "🟡 APPLIED"
	INTERVIEWING JobStatus = "🟠 INTERVIEWING"
	REJECTED     JobStatus = "🔴 REJECTED"
	SUCCEED      JobStatus = "🟢 SUCCEED"
)

var statusOrderSlice = []JobStatus{
	TO_APPLY,
	APPLIED,
	INTERVIEWING,
	SUCCEED,
	REJECTED,
}

func (status JobStatus) GetNextStatus() JobStatus {
	idx := status.findIndex()

	if idx+1 < len(statusOrderSlice) {
		return statusOrderSlice[idx+1]
	} else {
		return statusOrderSlice[0]
	}
}

func (status JobStatus) GetPrevStatus() JobStatus {
	idx := status.findIndex()

	if idx > 0 {
		return statusOrderSlice[idx-1]
	} else {
		return statusOrderSlice[len(statusOrderSlice)-1]
	}
}

func (status JobStatus) findIndex() int {
	for index, value := range statusOrderSlice {
		if value == status {
			return index
		}
	}
	return -1
}
