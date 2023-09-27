package snake

import (
	"context"
	"reflect"
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

func ResolveBindingsFromProvider(ctx context.Context, provider BindingResolver, rf reflect.Value) (context.Context, error) {

	// get the type of the first argument
	for i := 0; i < rf.Type().NumIn(); i++ {
		pt := rf.Type().In(i)
		if pt.Kind() == reflect.Ptr {
			pt = pt.Elem()
		}

		k := reflect.New(pt).Interface()
		p, err := provider.ResolveBinding(k)
		if err != nil {
			continue
		}
		if p != nil {
			ctx = Bind(ctx, k, p)
		}

	}

	return ctx, nil
}
