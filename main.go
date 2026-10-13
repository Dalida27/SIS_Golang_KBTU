package main

import (
	"fmt"

	"sis-project/task"
)

func main() {
	manager := task.NewManager(task.ScoreProcessor{PassScore: 50})

	fmt.Println("=== 1. Create a task ===")
	first, err := manager.CreateTask([]int{85, 42, 100, 67, 50})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("created task", first.ID, "with status", first.Status)

	fmt.Println()
	fmt.Println("=== 2. Process the task ===")
	err = manager.ProcessTask(first.ID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	printTask(first)

	fmt.Println()
	fmt.Println("=== 3. Invalid input is rejected ===")
	_, err = manager.CreateTask([]int{})
	fmt.Println("empty list ->", err)
	_, err = manager.CreateTask([]int{90, 150})
	fmt.Println("score 150  ->", err)

	fmt.Println()
	fmt.Println("=== 4. Invalid operations are rejected ===")
	err = manager.ProcessTask(first.ID)
	fmt.Println("process task 1 again ->", err)
	_, err = manager.GetTask(999)
	fmt.Println("get task 999         ->", err)

	fmt.Println()
	fmt.Println("=== 5. A second, independent task ===")
	second, err := manager.CreateTask([]int{30, 45, 60})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = manager.ProcessTask(second.ID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	printTask(second)

	fmt.Println()
	fmt.Println("=== 6. Every task keeps its own result ===")
	for id := 1; id <= 2; id++ {
		t, err := manager.GetTask(id)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Printf("task %d: status=%s average=%.2f passed=%d\n", t.ID, t.Status, t.Result.Average, t.Result.Passed)
	}
}

func printTask(t *task.Task) {
	fmt.Printf("task %d\n", t.ID)
	fmt.Printf("  scores:  %v\n", t.Scores)
	fmt.Printf("  status:  %s\n", t.Status)
	fmt.Printf("  count:   %d\n", t.Result.Count)
	fmt.Printf("  min:     %d\n", t.Result.Min)
	fmt.Printf("  max:     %d\n", t.Result.Max)
	fmt.Printf("  average: %.2f\n", t.Result.Average)
	fmt.Printf("  passed:  %d\n", t.Result.Passed)
}
