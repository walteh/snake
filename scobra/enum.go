package scobra

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"
)

var (
	_ pflag.Value = &wrappedEnum[string]{}
)

func NewWrappedEnum[I ~string | ~int](def I, curr *I, values ...I) (*wrappedEnum[I], error) {
	strt := &wrappedEnum[I]{values: values, current: curr}
	err := strt.set(def)
	if err != nil {
		return nil, err
	}
	return strt, nil
}

type wrappedEnum[I ~string | ~int] struct {
	current *I
	values  []I
}

// String is used both by fmt.Print and by Cobra in help text
func (e *wrappedEnum[I]) String() string {
	if e.current == nil {
		return ""
	}
	return fmt.Sprintf("%v", *e.current)
}

func (e *wrappedEnum[I]) convertFromString(v string) (I, error) {
	var vt I

	switch any(e.current).(type) {
	case *string:
		return any(v).(I), nil
	case *int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return vt, err
		}
		return any(i).(I), nil
	default:
		return vt, errors.Errorf("invalid type %T", e.current)
	}
}

func (e *wrappedEnum[I]) set(vt I) error {
	if slices.Contains(e.values, vt) {
		*e.current = vt
		return nil
	}
	return errors.Errorf("invalid value %q, expected one of %s", vt, strings.Join(e.ValuesStringSlice(), ", "))
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *wrappedEnum[I]) Set(v string) error {
	vt, err := e.convertFromString(v)
	if err != nil {
		return err
	}
	return e.set(vt)
}

// Type is only used in help text
func (e *wrappedEnum[I]) Type() string {
	if _, ok := any("").(I); ok {
		return "string"
	}
	if _, ok := any(0).(I); ok {
		return "int"
	}
	return "unknown"
}

func (e *wrappedEnum[I]) ValuesStringSlice() []string {
	wrk := make([]string, len(e.values))
	for i, v := range e.values {
		wrk[i] = fmt.Sprint(v)
	}
	return wrk
}

func (e *wrappedEnum[I]) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return e.ValuesStringSlice(), cobra.ShellCompDirectiveDefault
}

func (e *wrappedEnum[I]) Assign(cmd *cobra.Command, key string, descritpion string) error {
	cmd.Flags().Var(e, key, descritpion)
	err := cmd.RegisterFlagCompletionFunc(key, e.CompletionFunc)
	if err != nil {
		return err
	}
	return nil
}
