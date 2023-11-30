package scobra

import (
	"github.com/go-faster/errors"
	"github.com/manifoldco/promptui"
	"github.com/walteh/snake"
)

var _ snake.EnumResolverFunc = (*CobraSnake)(nil).ResolveEnum

func (me *CobraSnake) ResolveEnum(typ string, opts []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select " + typ,
		Items: opts,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	if result == "" {
		return "", errors.Errorf("invalid %q", typ)
	}

	return result, nil

}
