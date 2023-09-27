package snake

import (
	"context"
	"reflect"

	"github.com/spf13/cobra"
)

type BindingResolver interface {
	ResolveBinding(*cobra.Command, any) (any, error)
}

type bindingResolverKeyT struct {
}

func SetBindingResolver(ctx context.Context, provider BindingResolver) context.Context {
	return context.WithValue(ctx, &bindingResolverKeyT{}, provider)
}

func GetBindingResolver(ctx context.Context) BindingResolver {
	p, ok := ctx.Value(&bindingResolverKeyT{}).(BindingResolver)
	if ok {
		return p
	}
	return nil
}

func ResolveBindingsFromProvider(cmd *cobra.Command, rf reflect.Value, providers ...BindingResolver) error {

	ctx := cmd.Context()

	// get the type of the first argument
	for i := 0; i < rf.Type().NumIn(); i++ {
		pt := rf.Type().In(i)
		if pt.Kind() == reflect.Ptr {
			pt = pt.Elem()
		}

		k := reflect.New(pt).Interface()
		for _, provider := range providers {
			p, err := provider.ResolveBinding(cmd, k)
			if err != nil {
				continue
			}
			if p != nil {
				ctx = Bind(ctx, k, p)
				break
			}
		}
	}

	cmd.SetContext(ctx)

	return nil
}

type ResolverFunc[I any] func(*cobra.Command) (I, error)

type RawBindingResolver map[reflect.Type]ResolverFunc[any]

func (r RawBindingResolver) ResolveBinding(cmd *cobra.Command, key any) (any, error) {
	if f1, ok := r[reflect.TypeOf(key)]; ok {
		return f1(cmd)
	} else if f2, ok := r[reflect.TypeOf(key).Elem()]; ok {
		return f2(cmd)
	}
	return nil, nil
}

type dynamicBindingResolverKeyT struct {
}

func setDynamicBindingResolver(ctx context.Context, provider RawBindingResolver) context.Context {
	return context.WithValue(ctx, &dynamicBindingResolverKeyT{}, provider)
}

func getDynamicBindingResolver(ctx context.Context) RawBindingResolver {
	p, ok := ctx.Value(&dynamicBindingResolverKeyT{}).(RawBindingResolver)
	if ok {
		return p
	}
	return nil
}

func RegisterBindingResolver[I any](cmd *cobra.Command, resolver ResolverFunc[*I]) {
	// check if we have a dynamic binding resolver available
	ctx := cmd.Context()
	dy := getDynamicBindingResolver(ctx)
	if dy == nil {
		dy = RawBindingResolver{}
	}

	elm := reflect.TypeOf((*I)(nil)).Elem()

	dy[elm] = func(cmd *cobra.Command) (any, error) {
		return resolver(cmd)
	}

	ctx = setDynamicBindingResolver(ctx, dy)

	cmd.SetContext(ctx)
}
