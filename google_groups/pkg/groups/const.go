package groups

type GroupRole string

const (
	ManagerRole GroupRole = "MANAGER"
	OwnerRole GroupRole = "OWNER"
	MemberRole GroupRole = "MEMBER"
)
