package repository

func CheckTenderStatus(tenderStatus string) bool {
	return tenderStatus == "CREATED" || tenderStatus == "CLOSED" || tenderStatus == "PUBLISHED"
}
