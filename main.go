package main

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/jomei/notionapi"
	"github.com/joshrosso/nexp/export"
)

const (
	notionToken        = "NOTION_TOKEN"
	taskDBID           = "NOTION_TASK_DB_ID"
	propertyNameDate   = "Date"
	propertyNameStatus = "Status"
	propertyNameWork   = "Work"
	propertyNameTitle  = "Name"
	statusDone         = "done"
	templateFileName   = "slack-status.md.tmpl"
)

type Report struct {
	TodayTasks     []Task
	YesterdayTasks []Task
}

type Task struct {
	Title string
}

func main() {
	ctx := context.Background()
	client := notionapi.NewClient(notionapi.Token(os.Getenv(notionToken)))
	now := time.Now()
	r := Report{
		TodayTasks:     []Task{},
		YesterdayTasks: []Task{},
	}
	r.TodayTasks = getTodaysTasks(ctx, client, now)
	r.YesterdayTasks = getYesterdaysTasks(ctx, client, now)

	// Read the template file
	templateFile := templateFileName
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		panic(fmt.Sprintf("Error reading template file:", err))
	}

	// Render the template with the data
	err = tmpl.Execute(os.Stdout, &r)
	if err != nil {
		panic(fmt.Sprintf("Error rendering template:", err))
	}
}

func getYesterdaysTasks(ctx context.Context, client *notionapi.Client, now time.Time) []Task {
	// make now yesterday
	now = now.Add(-24 * time.Hour)
	// if sunday, bring in completed tasks from Friday
	if now.Weekday() == time.Sunday {
		now = now.Add(-48 * time.Hour)
	}
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)
	nTodayStart := notionapi.Date(todayStart)
	nTodayEnd := notionapi.Date(todayEnd)

	resp, err := client.Database.Query(ctx, notionapi.DatabaseID(os.Getenv(taskDBID)), &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: propertyNameDate,
				Date: &notionapi.DateFilterCondition{
					OnOrAfter: &nTodayStart,
				},
			},
			notionapi.PropertyFilter{
				Property: propertyNameDate,
				Date: &notionapi.DateFilterCondition{
					Before: &nTodayEnd,
				},
			},
			notionapi.PropertyFilter{
				Property: propertyNameStatus,
				Select: &notionapi.SelectFilterCondition{
					Equals: statusDone,
				},
			},
			notionapi.PropertyFilter{
				Property: propertyNameWork,
				Checkbox: &notionapi.CheckboxFilterCondition{
					Equals: true,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	r := export.MDRenderer{}
	tasks := []Task{}
	for _, task := range resp.Results {
		taskTitle := r.RenderText(task.Properties[propertyNameTitle].(*notionapi.TitleProperty).Title)
		tasks = append(tasks, Task{Title: taskTitle})
	}
	return tasks
}

func getTodaysTasks(ctx context.Context, client *notionapi.Client, now time.Time) []Task {
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)
	nTodayStart := notionapi.Date(todayStart)
	nTodayEnd := notionapi.Date(todayEnd)

	resp, err := client.Database.Query(ctx, notionapi.DatabaseID(os.Getenv(taskDBID)), &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: propertyNameDate,
				Date: &notionapi.DateFilterCondition{
					OnOrAfter: &nTodayStart,
				},
			},
			notionapi.PropertyFilter{
				Property: propertyNameDate,
				Date: &notionapi.DateFilterCondition{
					Before: &nTodayEnd,
				},
			},
			notionapi.PropertyFilter{
				Property: propertyNameWork,
				Checkbox: &notionapi.CheckboxFilterCondition{
					Equals: true,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	r := export.MDRenderer{}
	tasks := []Task{}
	for _, task := range resp.Results {
		taskTitle := r.RenderText(task.Properties[propertyNameTitle].(*notionapi.TitleProperty).Title)
		tasks = append(tasks, Task{Title: taskTitle})
	}
	return tasks
}
