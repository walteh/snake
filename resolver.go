package snake

import (
	"context"
	"reflect"

	"github.com/spf13/cobra"
)

type BindingResolver interface {
	ResolveBinding(any) (any, error)
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

func ResolveBindingsFromProvider(ctx context.Context, rf reflect.Value, providers ...BindingResolver) (context.Context, error) {

	// get the type of the first argument
	for i := 0; i < rf.Type().NumIn(); i++ {
		pt := rf.Type().In(i)
		if pt.Kind() == reflect.Ptr {
			pt = pt.Elem()
		}

		k := reflect.New(pt).Interface()
		for _, provider := range providers {
			p, err := provider.ResolveBinding(k)
			if err != nil {
				continue
			}
			if p != nil {
				ctx = Bind(ctx, k, p)
				break
			}
		}

	}

	return ctx, nil
}

type RawBindingResolver map[reflect.Type]func() (any, error)

func (r RawBindingResolver) ResolveBinding(key any) (any, error) {
	if f1, ok := r[reflect.TypeOf(key)]; ok {
		return f1()
	} else if f2, ok := r[reflect.TypeOf(key).Elem()]; ok {
		return f2()
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

func RegisterBindingResolver[I any](cmd *cobra.Command, resolver func() (*I, error)) {
	// check if we have a dynamic binding resolver available
	ctx := cmd.Context()
	dy := getDynamicBindingResolver(ctx)
	if dy == nil {
		dy = RawBindingResolver{}
	}

	elm := reflect.TypeOf((*I)(nil)).Elem()

	dy[elm] = func() (any, error) {
		return resolver()
	}

	ctx = setDynamicBindingResolver(ctx, dy)

	cmd.SetContext(ctx)
}
