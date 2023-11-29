package scobra

import "github.com/manifoldco/promptui"

func (me *wrappedEnum[A]) Run() (A, error) {

	if me.current == nil || *me.current == "select" {
		prompt := promptui.Select{
			Label: "Select ",
			Items: me.values,
		}

		_, result, err := prompt.Run()

		if err != nil {
			return A(""), err
		}

		return A(result), nil
	}

	return *me.current, nil
}
