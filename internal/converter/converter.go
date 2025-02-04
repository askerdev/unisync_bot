package converter

import (
	"strconv"
	"strings"
	"time"

	"github.com/askerdev/unisync_bot/internal/domain"
	"github.com/askerdev/unisync_bot/internal/mospolytech"
	"github.com/askerdev/unisync_bot/internal/tmpl"
)

func tasksFromLecture(
	chatID string,
	weekDay int,
	lectureNumber int,
	lecture *mospolytech.Lecture,
) ([]*domain.Task, error) {
	result := []*domain.Task{}

	start, err := time.Parse("2006-01-02", lecture.DateFrom)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", lecture.DateTo)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for curr := start; curr.Before(end); curr = curr.AddDate(0, 0, 1) {
		if int(curr.Weekday()) != weekDay {
			continue
		}
		t := &domain.Task{ChatID: chatID}

		class := []string{}
		for _, aud := range lecture.Auditories {
			class = append(class, aud.Title)
		}

		timeAt := time.Date(
			curr.Year(),
			curr.Month(),
			curr.Day(),
			domain.LectureHourMinute[lectureNumber][0],
			domain.LectureHourMinute[lectureNumber][1],
			0, 0,
			now.Location(),
		)
		t.TimeAt = timeAt.Unix()

		if t.TimeAt < now.Unix() {
			continue
		}

		params := &tmpl.MessageParams{
			Type:     lecture.Type,
			Subject:  lecture.Subject,
			TimeAt:   timeAt.Format("15:04 2006-01-02"),
			Teacher:  lecture.Teacher,
			Location: lecture.Location,
			Class:    strings.Join(class, ", "),
		}

		if lecture.Link != nil {
			params.Link = *lecture.Link
		}

		t.Text = tmpl.Message(params)

		result = append(result, t)
	}

	return result, nil
}

func TasksFromSchedule(
	chatID string,
	group string,
	sch *mospolytech.SemesterSchedule,
) ([]*domain.Task, error) {
	tasks := []*domain.Task{}
	for weekDayStr, day := range sch.Content[group].Grid {
		weekDay, err := strconv.Atoi(weekDayStr)
		if err != nil {
			return nil, err
		}
		for lectureNumberStr, ll := range day {
			lectureNumber, err := strconv.Atoi(
				lectureNumberStr,
			)
			if err != nil {
				return nil, err
			}
			for _, l := range ll {
				tasksFromLecture, err := tasksFromLecture(
					chatID,
					weekDay,
					lectureNumber,
					l,
				)
				if err != nil {
					return nil, err
				}
				tasks = append(
					tasks,
					tasksFromLecture...,
				)
			}
		}
	}
	return tasks, nil
}
