package scobra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake/sbind"
)

var (
	_ pflag.Value = &wrappedEnum{}
)

// func NewWrappedEnum(def I, curr *I, values ...I) (*wrappedEnum[I], error) {
// 	strt := &wrappedEnum[I]{values: values, current: curr}
// 	err := strt.set(def)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return strt, nil
// }

func NewWrappedEnum(opt sbind.EnumOption) *wrappedEnum {
	strt := &wrappedEnum{internal: opt}
	return strt
}

type wrappedEnum struct {
	internal sbind.EnumOption
}

// // String is used both by fmt.Print and by Cobra in help text
// func (e *wrappedEnum[I]) String() string {
// 	if e.current == nil {
// 		return ""
// 	}
// 	return fmt.Sprintf("%v", *e.current)
// }

// func (e *wrappedEnum[I]) convertFromString(v string) (I, error) {
// 	var vt I

// 	switch any(e.current).(type) {
// 	case *string:
// 		return any(v).(I), nil
// 	case *int:
// 		i, err := strconv.Atoi(v)
// 		if err != nil {
// 			return vt, err
// 		}
// 		return any(i).(I), nil
// 	default:
// 		return vt, errors.Errorf("invalid type %T", e.current)
// 	}
// }

// func (e *wrappedEnum[I]) set(vt I) error {
// 	if slices.Contains(e.internal.Options(), string(vt)) {
// 		*e.current = vt
// 		return nil
// 	}
// 	return errors.Errorf("invalid value %q, expected one of %s", vt, strings.Join(e.ValuesStringSlice(), ", "))
// }

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *wrappedEnum) Set(v string) error {
	return e.internal.SetCurrent(v)
}

func (e *wrappedEnum) String() string {
	if e.internal.CurrentPtr() == nil {
		return ""
	}
	return *e.internal.CurrentPtr()
}

// Type is only used in help text
func (e *wrappedEnum) Type() string {
	return "string"
}

func (e *wrappedEnum) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return e.internal.Options(), cobra.ShellCompDirectiveDefault
}

func (e *wrappedEnum) Assign(cmd *cobra.Command, key string, descritpion string) error {
	cmd.Flags().Var(e, key, descritpion)
	err := cmd.RegisterFlagCompletionFunc(key, e.CompletionFunc)
	if err != nil {
		return err
	}
	return nil
}
