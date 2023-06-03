package main

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/pothulapati/tailscale-talk/client"
	"github.com/pothulapati/tailscale-talk/client/todos"
	"github.com/pothulapati/tailscale-talk/pkg/models"
)

type checklistItem struct {
	todo       *models.Item
	widgetBool widget.Bool
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

	// sort payload by description
	list := todoResponse.GetPayload()
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[i].ID > list[j].ID {
				list[i], list[j] = list[j], list[i]
			}
		}
	}

	ui.TodoItems = make([]checklistItem, 0)
	for _, line := range list {
		ui.TodoItems = append(ui.TodoItems, checklistItem{
			todo: line,
			widgetBool: widget.Bool{
				Value: line.Completed,
			},
		})
	}

	checklistLayout := layout.List{
		Axis: layout.Vertical,
	}

	// Layout the list of checklist items.
	return checklistLayout.Layout(gtx, len(ui.TodoItems), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.CheckBox(ui.theme, &ui.TodoItems[i].widgetBool, *ui.TodoItems[i].todo.Description).Layout(gtx)
			}),
		)
	})
}
