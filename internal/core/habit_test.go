package core

import (
	"slices"
	"testing"
	"time"
)

func TestIsDue_DailyInterval1(t *testing.T) {
	habit, err := NewHabit("Читать книгу", "FREQ=DAILY;INTERVAL=1")
	if err != nil {
		t.Fatalf("не ожидали ошибку при создании привычки: %v", err)
	}

	due, err := habit.IsDue(time.Now())
	if err != nil {
		t.Fatalf("не ожидали ошибку в IsDue: %v", err)
	}

	if !due {
		t.Errorf("ожидали true для сегодняшнего дня, получили false")
	}
}

func TestIsDue_DailyInterval2(t *testing.T) {
	habit, err := NewHabit("Бег", "FREQ=DAILY;INTERVAL=2")
	if err != nil {
		t.Fatalf("не ожидали ошибку при создании привычки: %v", err)
	}

	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	dayAfterTomorrow := today.AddDate(0, 0, 2)

	// Сегодня -> true
	dueToday, err := habit.IsDue(today)
	if err != nil {
		t.Fatalf("ошибка при проверке сегодняшнего дня: %v", err)
	}
	if !dueToday {
		t.Errorf("на сегодня ожидали true, получили false")
	}

	// Завтра -> false
	dueTomorrow, err := habit.IsDue(tomorrow)
	if err != nil {
		t.Fatalf("ошибка при проверке завтрашнего дня: %v", err)
	}
	if dueTomorrow {
		t.Errorf("на завтра ожидали false, получили true")
	}

	// Послезавтра -> true
	dueDayAfter, err := habit.IsDue(dayAfterTomorrow)
	if err != nil {
		t.Fatalf("ошибка при проверке послезавтра: %v", err)
	}
	if !dueDayAfter {
		t.Errorf("на послезавтра ожидали true, получили false")
	}
}

func TestHabitDate(t *testing.T) {
	newDayBegins := 4

	expectedTime := time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)
	timeNormalized := HabitDate(time.Date(2026, 11, 26, 3, 0, 0, 0, time.UTC), newDayBegins)
	if !timeNormalized.Equal(expectedTime) {
		t.Errorf("ожидали %v, получили %v", expectedTime, timeNormalized)
	}

	expectedTime = time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC)
	timeNormalized = HabitDate(time.Date(2026, 11, 26, 4, 0, 0, 0, time.UTC), newDayBegins)
	if !timeNormalized.Equal(expectedTime) {
		t.Errorf("ожидали %v, получили %v", expectedTime, timeNormalized)
	}

	expectedTime = time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC)
	timeNormalized = HabitDate(time.Date(2026, 11, 26, 5, 0, 0, 0, time.UTC), newDayBegins)
	if !timeNormalized.Equal(expectedTime) {
		t.Errorf("ожидали %v, получили %v", expectedTime, timeNormalized)
	}

}

func TestIsDue_Weekly(t *testing.T) {
	habit := Habit{
		Title:     "Weekly Habit, MO,WE,FR",
		RRule:     "FREQ=WEEKLY;BYDAY=MO,WE,FR",
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	expectedDays := map[time.Weekday]bool{
		time.Monday:    true,
		time.Tuesday:   false,
		time.Wednesday: true,
		time.Thursday:  false,
		time.Friday:    true,
		time.Saturday:  false,
		time.Sunday:    false,
	}

	for i := 0; i < 7; i++ {
		targetDate := habit.StartDate.AddDate(0, 0, i)
		day := targetDate.Weekday()
		expected := expectedDays[day]

		due, err := habit.IsDue(targetDate)

		if err != nil {
			t.Fatalf("ошибка при проверке дня %v: %v", day, err)
		}
		if due != expected {
			t.Errorf("на день %v ожидали %v, получили %v", day, expected, due)
		}

	}
}

// TestNewHabit проверяет конструктор NewHabit:
// корректное заполнение полей при валидном правиле и возврат ошибки при невалидном.
func TestNewHabit(t *testing.T) {
	// t.Run позволяет разбить один тест на логические подтесты.
	// В отчёте go test -v каждый подтест будет виден отдельной строкой,
	// и при падении одного сразу понятно, какой именно сценарий сломался.
	t.Run("валидная привычка создаётся без ошибок", func(t *testing.T) {
		title := "Зарядка"
		rrule := "FREQ=DAILY;INTERVAL=1"

		habit, err := NewHabit(title, rrule)
		// Если создание вернуло ошибку, проверять поля habit бессмысленно (там nil),
		// поэтому сразу останавливаем этот подтест через t.Fatalf.
		if err != nil {
			t.Fatalf("не ожидали ошибку при создании привычки, получили: %v", err)
		}
		if habit == nil {
			t.Fatalf("ожидали непустой указатель на Habit, получили nil")
		}

		// Дальше идут независимые проверки полей — используем t.Errorf,
		// чтобы увидеть сразу все несовпадения, если сломалось несколько полей.
		if habit.Title != title {
			t.Errorf("Title: ожидали %q, получили %q", title, habit.Title)
		}
		if habit.RRule != rrule {
			t.Errorf("RRule: ожидали %q, получили %q", rrule, habit.RRule)
		}
		if habit.IsArchived {
			t.Errorf("новая привычка не должна быть архивной (IsArchived == false)")
		}
		if habit.UUID == "" {
			t.Errorf("UUID должен быть сгенерирован и не быть пустым")
		}
		if habit.StartDate.IsZero() {
			t.Errorf("StartDate должна быть заполнена (не нулевое время)")
		}
	})

	t.Run("невалидный RRule возвращает ошибку", func(t *testing.T) {
		// Передаём заведомо сломанное правило (неизвестный ключ и без знака '=').
		habit, err := NewHabit("Сломанная", "BROKEN_RRULE")

		// Здесь мы, наоборот, ОЖИДАЕМ ошибку: если err == nil, это баг.
		if err == nil {
			t.Fatalf("ожидали ошибку для невалидного RRule, но получили nil")
		}
		if habit != nil {
			t.Errorf("при ошибке конструктор должен возвращать nil вместо %+v", habit)
		}
	})
}

// TestIsDue_EdgeCases — табличный тест (table-driven test) для граничных случаев IsDue.
// Вместо копирования похожего кода в 5 отдельных функций мы описываем сценарии
// в виде слайса структур и прогоняем их в одном цикле через t.Run.
func TestIsDue_EdgeCases(t *testing.T) {
	// Фиксируем опорную дату старта: 25 ноября 2026 года, 18:00 UTC (среда).
	startDate := time.Date(2026, 11, 25, 18, 0, 0, 0, time.UTC)

	// Каждый элемент слайса — отдельный тестовый сценарий (входные данные + ожидания).
	tests := []struct {
		name       string    // название подтеста (будет отображаться в выводе go test -v)
		rrule      string    // правило повторения
		isArchived bool      // флаг архивации
		targetDate time.Time // дата, которую проверяем
		wantDue    bool      // ожидаемый булевый результат IsDue
		wantErr    bool      // ожидаем ли ошибку (true) или nil (false)
	}{
		{
			name:       "архивная привычка никогда не активна",
			rrule:      "FREQ=DAILY;INTERVAL=1",
			isArchived: true,
			// Дата полностью подходит под DAILY, но из-за IsArchived должно быть false.
			targetDate: startDate.AddDate(0, 0, 1),
			wantDue:    false,
			wantErr:    false,
		},
		{
			name:       "целевая дата раньше StartDate",
			rrule:      "FREQ=DAILY;INTERVAL=1",
			isArchived: false,
			// Привычка ещё не началась (проверяем вчерашний день относительно старта).
			targetDate: startDate.AddDate(0, 0, -1),
			wantDue:    false,
			wantErr:    false,
		},
		{
			name:       "день старта, но время раньше часа создания",
			rrule:      "FREQ=DAILY;INTERVAL=1",
			isArchived: false,
			// StartDate в 18:00, а проверяем в 09:00 того же дня.
			// Благодаря dateOnly() время суток должно игнорироваться -> true.
			targetDate: time.Date(2026, 11, 25, 9, 0, 0, 0, time.UTC),
			wantDue:    true,
			wantErr:    false,
		},
		{
			name:       "WEEKLY без указания BYDAY возвращает false",
			rrule:      "FREQ=WEEKLY",
			isArchived: false,
			targetDate: startDate,
			wantDue:    false,
			wantErr:    false,
		},
		{
			name:       "неподдерживаемая частота (MONTHLY) возвращает ошибку",
			rrule:      "FREQ=MONTHLY;INTERVAL=1",
			isArchived: false,
			targetDate: startDate,
			wantDue:    false,
			wantErr:    true,
		},
		{
			name:       "некорректный RRule в структуре возвращает ошибку",
			rrule:      "BAD_RULE",
			isArchived: false,
			targetDate: startDate,
			wantDue:    false,
			wantErr:    true,
		},
	}

	// Проходим по всем сценариям из таблицы.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			habit := Habit{
				Title:      "Тестовая привычка",
				RRule:      tc.rrule,
				StartDate:  startDate,
				IsArchived: tc.isArchived,
			}

			gotDue, err := habit.IsDue(tc.targetDate)

			// (err != nil) превращает наличие ошибки в bool и позволяет
			// одной строкой проверить оба случая: ждали ошибку или не ждали.
			if (err != nil) != tc.wantErr {
				t.Fatalf("IsDue() ошибка = %v, ожидали ошибку: %v", err, tc.wantErr)
			}
			if gotDue != tc.wantDue {
				t.Errorf("IsDue() = %v, ожидали %v", gotDue, tc.wantDue)
			}
		})
	}
}

// TestParseRRule проверяет внутренний парсер RRULE-строк.
// Так как тестовый файл находится в том же пакете (package core),
// у нас есть прямой доступ к неэкспортируемой функции parseRRule и полям recurrenceRule.
func TestParseRRule(t *testing.T) {
	t.Run("полный валидный RRULE разбирается корректно", func(t *testing.T) {
		// Включаем все поддерживаемые ключи и лишнюю ';' в конце (должна игнорироваться).
		input := "FREQ=WEEKLY;INTERVAL=2;UNTIL=20301231T235959Z;COUNT=10;BYDAY=MO,WE;"

		rule, err := parseRRule(input)
		if err != nil {
			t.Fatalf("не ожидали ошибку парсинга, получили: %v", err)
		}

		if rule.Freq != "WEEKLY" {
			t.Errorf("Freq: ожидали WEEKLY, получили %q", rule.Freq)
		}
		if rule.Interval != 2 {
			t.Errorf("Interval: ожидали 2, получили %d", rule.Interval)
		}
		if rule.count != 10 {
			t.Errorf("count: ожидали 10, получили %d", rule.count)
		}

		expectedUntil := time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)
		if !rule.until.Equal(expectedUntil) {
			t.Errorf("until: ожидали %v, получили %v", expectedUntil, rule.until)
		}

		// Сравниваем слайсы строк через стандартный пакет slices.
		expectedDays := []string{"MO", "WE"}
		if !slices.Equal(rule.byDays, expectedDays) {
			t.Errorf("byDays: ожидали %v, получили %v", expectedDays, rule.byDays)
		}
	})

	t.Run("интервал по умолчанию равен 1", func(t *testing.T) {
		rule, err := parseRRule("FREQ=DAILY")
		if err != nil {
			t.Fatalf("не ожидали ошибку парсинга, получили: %v", err)
		}
		if rule.Interval != 1 {
			t.Errorf("Interval по умолчанию должен быть 1, получили %d", rule.Interval)
		}
	})

	// Табличная проверка всех вариантов ошибочного ввода для parseRRule.
	invalidCases := []struct {
		name  string
		rrule string
	}{
		{
			name:  "часть без знака равно",
			rrule: "FREQ",
		},
		{
			name:  "нечисловой INTERVAL",
			rrule: "FREQ=DAILY;INTERVAL=abc",
		},
		{
			name:  "неверный формат даты UNTIL",
			rrule: "FREQ=DAILY;UNTIL=2030-12-31",
		},
		{
			name:  "нечисловой COUNT",
			rrule: "FREQ=DAILY;COUNT=много",
		},
		{
			name:  "неизвестный ключ",
			rrule: "FREQ=DAILY;UNKNOWN=1",
		},
	}

	for _, tc := range invalidCases {
		t.Run("ошибка при: "+tc.name, func(t *testing.T) {
			rule, err := parseRRule(tc.rrule)
			if err == nil {
				t.Errorf("для %q ожидали ошибку, но получили nil (rule=%+v)", tc.rrule, rule)
			}
		})
	}
}

// TestHabitDate_Boundaries дополняет TestHabitDate проверкой перехода через
// полночь, границу месяца/года и нулевого часа старта дня.
func TestHabitDate_Boundaries(t *testing.T) {
	tests := []struct {
		name         string
		input        time.Time
		dayStartHour int
		want         time.Time
	}{
		{
			name: "переход через границу года (1 января 02:00 при старте в 06:00 -> 31 декабря)",
			// AddDate(0, 0, -1) должен корректно переключить не только день, но и месяц с годом.
			input:        time.Date(2026, 1, 1, 2, 0, 0, 0, time.UTC),
			dayStartHour: 6,
			want:         time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "ровно полночь при старте дня в 06:00 относится к вчерашнему дню",
			// 0 < 6 — это всё ещё вчерашний день привычек.
			input:        time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC),
			dayStartHour: 6,
			want:         time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "старт дня в 00:00 никогда не сдвигает дату назад",
			// В полночь Hour() == 0, условие 0 < 0 ложно -> остаётся 26 ноября.
			input:        time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC),
			dayStartHour: 0,
			want:         time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HabitDate(tc.input, tc.dayStartHour)
			// Всегда сравниваем time.Time через .Equal(), а не через ==.
			if !got.Equal(tc.want) {
				t.Errorf("HabitDate(%v, %d) = %v, ожидали %v", tc.input, tc.dayStartHour, got, tc.want)
			}
		})
	}
}
