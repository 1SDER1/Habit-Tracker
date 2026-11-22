package core

/*
1. Добавтиь правила RRULE для привычек в функции IsHabitDue
*/
import "time"

type Habit struct {
	UUID       string
	Title      string
	RRule      string
	StartDate  time.Time
	IsArchived bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type Entry struct {
	UUID        string
	HabitUUID   string
	Date        string
	CompletedAt time.Time
	UpdatedAt   time.Time
}

func IsHabitDue(h Habit, targetDate time.Time) (bool, error) {
	if h.IsArchived {
		return false, nil
	}

	if truncateToDay(targetDate).Before(truncateToDay(h.StartDate)) {
		return false, nil
	}
	// Правила RRULE добавть
	return true, nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func NormalizeToHabitDate(t time.Time, dayStartHour int) time.Time {
	// Входит ли в привычку день, если привычка начинается в 6 утра, а сейчас 5 утра
	// , то привычка считается за вчерашний день
	if t.Hour() < dayStartHour {
		t = t.AddDate(0, 0, -1)
	}
	return truncateToDay(t)
}
