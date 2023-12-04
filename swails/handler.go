package swails

import (
	"os"

	"github.com/go-faster/errors"
	"github.com/walteh/snake"
)

type WailsHTMLResponse struct {
	Default     string     `json:"default"`
	Text        string     `json:"text"`
	JSON        any        `json:"json"`
	Table       [][]string `json:"table"`
	TableStyles [][]string `json:"table_styles"`
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
	Name   string          `json:"name"`
	Type   snake.InputType `json:"type"`
	Value  any             `json:"value"`
	Shared bool            `json:"shared"`
}

func (me *WailsSnake) Inputs() ([]*WailsInput, error) {

	var inputs []*WailsInput
	for _, input := range me.inputs {
		if input.Shared() {
			inputs = append(inputs, &WailsInput{
				Name:   input.Name(),
				Type:   input.Type(),
				Value:  input.Ptr(),
				Shared: input.Shared(),
			})
		}
	}

	return inputs, nil
}

func (me *WailsSnake) InputsFor(name *WailsCommand) ([]*WailsInput, error) {
	cmd := me.snake.Resolve(name.Name)

	snk, err := snake.InputsFor(me.snake.Resolve(name.Name), me.snake.Enums()...)
	if err != nil {
		return nil, err
	}

	var wail SWails
	if typ, ok := cmd.(snake.TypedResolver[SWails]); ok {
		wail = typ.TypedRef()
	} else {
		return nil, errors.Errorf("command %q is not a wails command", name.Name)
	}

	var inputs []*WailsInput
	for _, input := range snk {
		inp, err := me.CurrentInput(wail, input)
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

func (me *WailsSnake) CurrentInput(cmd SWails, input snake.Input) (*WailsInput, error) {

	curr := me.inputs[inputName(cmd, input)]

	return &WailsInput{
		Name:   inputName(cmd, input),
		Type:   curr.Type(),
		Value:  curr.Ptr(),
		Shared: curr.Shared(),
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

	inp := me.inputs[input.Name]

	return &WailsInput{
		Name:   input.Name,
		Type:   inp.Type(),
		Value:  inp.Ptr(),
		Shared: inp.Shared(),
	}, nil
}

type WailsCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Emoji       string `json:"emoji"`
}

func (me *WailsSnake) Commands() ([]*WailsCommand, error) {
	var cmds []*WailsCommand
	for _, cmd := range me.snake.Resolvers() {
		if cmdt, ok := cmd.(snake.TypedResolver[SWails]); ok {
			cmds = append(cmds, &WailsCommand{
				Name:        cmdt.TypedRef().Name(),
				Description: cmdt.TypedRef().Description(),
				Image:       cmdt.TypedRef().Image(),
				Emoji:       cmdt.TypedRef().Emoji(),
			})
		}
	}

	return cmds, nil
}
