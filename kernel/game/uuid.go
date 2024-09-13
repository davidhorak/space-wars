package game

var uuid int64 = 0

func NewUUID() int64 {
	uuid++
	return uuid
}
