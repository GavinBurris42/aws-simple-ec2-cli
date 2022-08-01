// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
package questionModel

import (
	"simple-ec2/pkg/cli"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

/*
SingleSelectList represents a question with a list of options from which a single option is chosen as the answer.
Options may be presented in a table based on initialized input.
*/
type Confirmation struct {
	lists      []SingleSelectList
	choice     string // The chosen option
	focusIndex int
	allowEdit  bool
	err        error // An error caught during the question
}

// InitializeModel initializes the model based on the passed in question input
func (c *Confirmation) InitializeModel(input *QuestionInput) {
	configList := SingleSelectList{}
	configList.InitializeModel(&QuestionInput{
		HeaderStrings:  []string{"Configurations", "Values"},
		QuestionString: "Please confirm if you would like to launch instance with following options:",
		Rows:           input.Rows,
		IndexedOptions: input.IndexedOptions,
	})
	configList.list.Select(-1)

	yesNoList := SingleSelectList{}
	yesNoList.InitializeModel(&QuestionInput{
		IndexedOptions: yesNoOptions,
		DefaultOption:  cli.ResponseNo,
		Rows:           CreateSingleLineRows(yesNoData),
	})
	c.lists = append(c.lists, configList, yesNoList)
	c.focusIndex = 1
}

// Init defines an optional command that can be run when the question is asked.
func (c *Confirmation) Init() tea.Cmd {
	return nil
}

/*
Update is called when a message is received. Handles user input to traverse through list and
select an answer.
*/
func (c *Confirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			c.err = exitError
			return c, tea.Quit
		case tea.KeyUp:
			if c.focusIndex > -len(c.lists[0].list.Items()) && (c.allowEdit || c.focusIndex > 0) {
				c.focusIndex--
			}
			if c.focusIndex == -1 {
				c.lists[0].list.Select(len(c.lists[0].list.Items()))
				c.lists[1].list.Select(c.focusIndex)
			}
		case tea.KeyDown:
			if c.focusIndex < len(c.lists[1].list.Items())-1 {
				c.focusIndex++
			}
			if c.focusIndex == 0 {
				c.lists[0].list.Select(-1)
				c.lists[1].list.Select(c.focusIndex)
				return c, nil
			}
		case tea.KeyEnter:
			if c.focusIndex < 0 {
				c.lists[0].selectItem()
				c.choice = c.lists[0].GetChoice()
				return c, tea.Quit
			} else {
				c.lists[1].selectItem()
				c.choice = c.lists[1].GetChoice()
				return c, tea.Quit
			}
		}

	case error:
		c.err = msg
		return c, tea.Quit
	}

	if c.focusIndex < 0 && c.allowEdit {
		c.lists[0].Update(msg)
	} else {
		c.lists[1].Update(msg)
	}
	return c, nil
}

// View renders the view for the question. The view is rendered after every update
func (c *Confirmation) View() string {
	b := strings.Builder{}
	b.WriteString(c.lists[0].View())
	b.WriteRune('\n')
	b.WriteString(c.lists[1].View())
	return b.String()
}

// GetChoice gets the selected choice
func (c *Confirmation) GetChoice() string { return c.choice }

// getError gets the error from the question if one arose
func (c *Confirmation) getError() error { return c.err }

func (c *Confirmation) SetAllowEdit(allowEdit bool) {
	c.allowEdit = allowEdit
}
