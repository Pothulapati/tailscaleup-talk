package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func getTailScaleServer() string {
	// List all the servers on a tag
	return "todo-server"
}

func (ui *UI) layoutOAuth(gtx layout.Context) layout.Dimensions {
	// Define the input fields.
	var email, password widget.Editor
	email.SingleLine, password.SingleLine = true, true
	email.Submit, password.Submit = true, true

	// Define the sign-in button.
	signin := material.Button(ui.theme, nil, "Sign in")

	// Layout the input fields and sign-in button.
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return material.Editor(ui.theme, &email, "Email").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Editor(ui.theme, &password, "Password").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return signin.Layout(gtx)
		}),
	)
}
