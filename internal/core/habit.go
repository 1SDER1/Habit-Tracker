package core

/*
1. Добавтиь правила RRULE для привычек в функции IsHabitDue
*/
import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

func NewHabit(title string, rrule string) (*Habit, error) {

	_, err := parseRRule(rrule)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	habit := Habit{
		UUID:       strconv.Itoa(int(time.Now().UnixNano())), //TODO: generate UUID
		Title:      title,
		RRule:      rrule,
		IsArchived: false,
		StartDate:  now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return &habit, nil
}

type recurrenceRule struct {
	Freq     string    //YEARLY MONTHLY WEEKLY DAILY HOURLY MINUTELY SECONDLY
	Interval int       //по умолчанию 1, при freq 2 и weekly, интервал 2 означает каждые 2 недели
	byDays   []string  //дни недели, в которые повторяется привычка, например "MO", "TU", "WE", "TH", "FR", "SA", "SU"
	until    time.Time // дата, до которой повторяется привычка
	count    int       // количество повторений

}

func parseRRule(rrule string) (*recurrenceRule, error) {
	//FREQ=YEARLY;INTERVAL=1;UNTIL=20301231T235959Z;COUNT=10;BYDAY=MO,WE;
	rule := &recurrenceRule{Interval: 1} // по умолчанию интервал равен 1
	parts := strings.Split(rrule, ";")
	for _, part := range parts {

		if part == "" {
			continue
		}

		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid rrule part: %s", part)
		}

		key, value := keyValue[0], keyValue[1]

		switch key {
		case "FREQ":
			rule.Freq = value

		case "INTERVAL":
			interval, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid INTERVAL value: %s", value)
			}
			rule.Interval = interval

		case "BYDAY":
			rule.byDays = strings.Split(value, ",")

		case "UNTIL":
			until, err := time.Parse("20060102T150405Z", value)
			if err != nil {
				return nil, fmt.Errorf("invalid UNTIL value: %s", value)
			}
			rule.until = until

		case "COUNT":
			count, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid COUNT value: %s", value)
			}
			rule.count = count

		default:
			return nil, fmt.Errorf("unknown rrule key: %s", key)
		}
	}
	return rule, nil
}

func (h Habit) IsDue(targetDate time.Time) (bool, error) {
	//Базовые проверки
	if h.IsArchived {
		return false, nil
	}

	if dateOnly(targetDate).Before(dateOnly(h.StartDate)) {
		return false, nil
	}

	//Проверка по правилам RRULE
	rule, err := parseRRule(h.RRule)
	if err != nil {
		return false, err
	}
	switch rule.Freq {
	case "DAILY":
		if rule.Interval == 1 {
			return true, nil
		}
		// Проверяем, является ли целевая дата кратной интервалу с момента начала привычки
		if int(dateOnly(targetDate).Sub(dateOnly(h.StartDate)).Hours()/24)%rule.Interval == 0 {
			return true, nil
		}
		return false, nil
	case "WEEKLY":
		if rule.contains(weekdayToICalendar(targetDate.Weekday())) {
			return true, nil
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported frequency: %s", rule.Freq)
	}
	//FREQ=YEARLY;INTERVAL=1;UNTIL=20301231T235959Z;COUNT=10;BYDAY=MO,WE;
	// Правила RRULE добавть
}

func (r *recurrenceRule) contains(day string) bool {
	for _, d := range r.byDays {
		if d == day {
			return true
		}
	}
	return false
}

func weekdayToICalendar(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "MO"
	case time.Tuesday:
		return "TU"
	case time.Wednesday:
		return "WE"
	case time.Thursday:
		return "TH"
	case time.Friday:
		return "FR"
	case time.Saturday:
		return "SA"
	case time.Sunday:
		return "SU"
	default:
		return ""
	}
}

func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func HabitDate(t time.Time, dayStartHour int) time.Time {
	// Входит ли в привычку день, если привычка начинается в 6 утра, а сейчас 5 утра
	// , то привычка считается за вчерашний день
	if t.Hour() < dayStartHour {
		t = t.AddDate(0, 0, -1)
	}
	return dateOnly(t)
}
