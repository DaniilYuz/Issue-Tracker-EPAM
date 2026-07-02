package utils

func UserKey(id string) string {
	return "users:" + id
}

func IssueKey(id string) string {
	return "issues:" + id
}

func ProjectKey(id string) string {
	return "projects:" + id
}
