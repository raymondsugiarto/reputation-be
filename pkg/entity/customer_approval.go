package entity

// ApprovalActionRequestDto is the request body for both
// POST /api/admin/approvals/:customerId/approve and .../reject.
//
// `remark` is required at the handler layer for the reject action
// (admin must justify rejection), but optional for approve. Service
// code enforces the per-action rules.
type ApprovalActionRequestDto struct {
	Remark string `json:"remark"`
}

// CustomerApprovalStatsDto powers the approval page's counter strip.
// Counters are platform-wide and do not require a date range — admins
// only need a current snapshot to triage the queue.
type CustomerApprovalStatsDto struct {
	PendingApprovals int64 `json:"pendingApprovals"`
	ApprovedToday    int64 `json:"approvedToday"`
	RejectedToday    int64 `json:"rejectedToday"`
	TotalApproved    int64 `json:"totalApproved"`
	TotalRejected    int64 `json:"totalRejected"`
}
