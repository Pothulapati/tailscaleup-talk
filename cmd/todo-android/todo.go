package main

import (
	"fmt"
	"sort"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/aws/smithy-go/time"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/pothulapati/tailscale-talk/client"
	"github.com/pothulapati/tailscale-talk/client/todos"
	"github.com/pothulapati/tailscale-talk/pkg/models"
)

type checklistItem struct {
	clickable widget.Clickable
	todo      *models.Item
}

type TodoUI struct {
	checklistItems []checklistItem
}

func (item *checklistItem) layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// Create a checkbox widget for the item.
	checkbox := material.CheckBox(theme, &widget.Bool{
		Value: item.todo.Completed,
	}, *item.todo.Description)

	// Detect clicks on the checkbox.
	for item.clickable.Clicked() {
		// Update the todo completion status.
		// Update the todo using the client.
		transport := httptransport.New("todo", "/", []string{"http"})
		transport.Producers["application/io.goswagger.examples.todo-list.v1+json"] = runtime.JSONProducer()
		transport.Consumers["application/io.goswagger.examples.todo-list.v1+json"] = runtime.JSONConsumer()

		client := client.New(transport, nil)

		_, err := client.Todos.UpdateOne(todos.NewUpdateOneParams().WithID(item.todo.ID).WithBody(&models.Item{
			Description: item.todo.Description,
			Completed:   item.todo.Completed,
		}))
		if err != nil {
			panic(err)
		}
	}

	// Layout the checkbox widget.
	return checkbox.Layout(gtx)
}

func (ui *UI) layoutTodo(gtx layout.Context) layout.Dimensions {
	transport := httptransport.New(getTailScaleServer(), "/", []string{"http"})
	transport.Producers["application/io.goswagger.examples.todo-list.v1+json"] = runtime.JSONProducer()
	transport.Consumers["application/io.goswagger.examples.todo-list.v1+json"] = runtime.JSONConsumer()

	client := client.New(transport, nil)

	todoResponse, err := client.Todos.FindTodos(todos.NewFindTodosParams())
	if err != nil {
		return layout.Center.Layout(gtx, material.Body1(ui.theme, fmt.Sprintf("Error fetching todos: %s", err.Error())).Layout)
	}

	// sort the todos by completion status & then by creation date
	sort.Slice(todoResponse.GetPayload(), func(i, j int) bool {
		if todoResponse.GetPayload()[i].Completed == todoResponse.GetPayload()[j].Completed {
			x, _ := time.ParseDateTime(todoResponse.GetPayload()[i].CreatedAt.String())
			y, _ := time.ParseDateTime(todoResponse.GetPayload()[j].CreatedAt.String())
			return x.Before(y)
		}
		return todoResponse.GetPayload()[i].Completed
	})

	var items []checklistItem
	for _, line := range todoResponse.GetPayload() {
		items = append(items, checklistItem{
			todo: line,
		})
	}

	checklistLayout := layout.List{
		Axis: layout.Vertical,
	}

	// Layout the list of checklist items.
	return checklistLayout.Layout(gtx, len(items), func(gtx layout.Context, index int) layout.Dimensions {
		item := items[index]
		return item.layout(gtx, ui.theme)
	})
}
