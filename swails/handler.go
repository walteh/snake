package swails

import (
	"os"

	"github.com/go-faster/errors"
	"github.com/walteh/snake"
)

type WailsHTMLResponse struct {
	HTML string `json:"html"`
}

// the handler is designed to be a wails binding that will automatically inject snake bindings
func (me *WailsSnake) Run(name *WailsCommand) (*WailsHTMLResponse, error) {
	outhand := NewOutputHandler(os.Stdout)

	err := snake.RunResolvingArguments(outhand, me.snake.Resolve, name.Name, me.binder)
	if err != nil {
		return nil, err
	}

	return outhand.output, nil
}

type WailsInput struct {
	Name  string          `json:"name"`
	Type  snake.InputType `json:"type"`
	Value any             `json:"value"`
}

func (me *WailsSnake) Inputs() ([]*WailsInput, error) {

	var inputs []*WailsInput
	for _, input := range me.inputs {
		inp, err := me.CurrentInput(input.Name())
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, inp)
	}

	return inputs, nil
}

func (me *WailsSnake) InputsFor(name *WailsCommand) ([]*WailsInput, error) {
	snk, err := snake.InputsFor(me.snake.Resolve(name.Name), me.snake.Enums()...)
	if err != nil {
		return nil, err
	}

	var inputs []*WailsInput
	for _, input := range snk {
		inp, err := me.CurrentInput(input.Name())
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, inp)
	}

	return inputs, nil
}

func (me *WailsSnake) OptionsForEnum(input *WailsInput) ([]string, error) {

	if input.Type != snake.StringEnumInputType {
		return nil, errors.Errorf("input %q is not an enum", input.Name)
	}

	snk := me.snake.Enums()

	var options []string
	for _, option := range snk {
		if option.Name() != input.Name {
			continue
		}
		options = append(options, option.Options()...)
	}

	return options, nil
}

func (me *WailsSnake) CurrentInput(name string) (*WailsInput, error) {

	curr := me.inputs[name]

	return &WailsInput{
		Name:  curr.Name(),
		Type:  curr.Type(),
		Value: curr.Ptr(),
	}, nil
}

func (me *WailsSnake) UpdateInput(input *WailsInput) (*WailsInput, error) {

	curr := me.inputs[input.Name]

	err := curr.SetValue(input.Value)
	if err != nil {
		return nil, errors.Errorf("unable to update input %q: %w", input.Name, err)
	}

	me.binder, err = snake.RefreshDependencies(curr, me.snake, me.binder)
	if err != nil {
		return nil, errors.Errorf("unable to update input %q: %w", input.Name, err)
	}

	return me.CurrentInput(input.Name)
}

type WailsCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (me *WailsSnake) Commands() ([]*WailsCommand, error) {
	var cmds []*WailsCommand
	for _, cmd := range me.snake.Resolvers() {
		if cmdt, ok := cmd.(snake.TypedResolver[SWails]); ok {
			cmds = append(cmds, &WailsCommand{
				Name:        cmdt.TypedRef().Name(),
				Description: cmdt.TypedRef().Description(),
			})
		}
	}

	return cmds, nil
}
