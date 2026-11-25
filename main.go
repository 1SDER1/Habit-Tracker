package main

import (
	"fmt"
	"habitcalendar/internal/core"
	"time"
)

func main() {
	habit, err := core.NewHabit("Пить воду", "FREQ=YEARLY;INTERVAL=1;UNTIL=20301231T235959Z;COUNT=10;BYDAY=MO,WE;")
	if err != nil {
		panic(err)
	}
	normalizedDate := core.HabitDate(time.Now(), 6)
	isDue, err := habit.IsDue(normalizedDate)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Is habit due: %t\n", isDue)
}
