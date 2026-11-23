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

func (h Habit) createHabit(title string, freq string) error {
	return nil
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
