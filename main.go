package main

import (
	"fmt"
	"habitcalendar/internal/core"
	"time"
)

func main() {
	habit := core.Habit{
		UUID:       "1",
		Title:      "Пить воду",
		RRule:      "FREQ=DAILY",
		StartDate:  time.Now(),
		IsArchived: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	normalizedDate := core.NormalizeToHabitDate(time.Now(), 6)
	isDue, err := core.IsHabitDue(habit, normalizedDate)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Is habit due: %t\n", isDue)
}
