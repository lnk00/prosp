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
