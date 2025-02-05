package domain

import "time"

type Task struct {
	ID     int
	ChatID string
	Text   string
	TimeAt time.Time
}

var LectureHourMinute = map[int][]int{
	1: {9, 0},
	2: {10, 40},
	3: {12, 20},
	4: {14, 30},
	5: {16, 10},
	6: {17, 50},
	7: {19, 30},
}
