package snake

import (
	"context"
	"reflect"
)

type BindingResolver interface {
	ResolveBinding(context.Context, any) (any, error)
}

type bindingResolverKeyT struct {
}

// func SetBindingResolver(ctx context.Context, provider BindingResolver) context.Context {
// 	return context.WithValue(ctx, &bindingResolverKeyT{}, provider)
// }

// func GetBindingResolver(ctx context.Context) BindingResolver {
// 	p, ok := ctx.Value(&bindingResolverKeyT{}).(BindingResolver)
// 	if ok {
// 		return p
// 	}
// 	return nil
// }

func ResolveBindingsFromProvider(ctx context.Context, rf reflect.Value, providers ...BindingResolver) (context.Context, error) {

	for i := 0; i < rf.Type().NumIn(); i++ {
		pt := rf.Type().In(i)
		if pt.Kind() == reflect.Ptr {
			pt = pt.Elem()
		}

		k := reflect.New(pt).Interface()
		for _, provider := range providers {
			p, err := provider.ResolveBinding(ctx, k)
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

type ResolverFunc[I any] func(ctx context.Context) (I, error)

type RawBindingResolver map[reflect.Type]ResolverFunc[any]

func (r RawBindingResolver) ResolveBinding(ctx context.Context, key any) (any, error) {
	if f1, ok := r[reflect.TypeOf(key)]; ok {
		return f1(ctx)
	} else if f2, ok := r[reflect.TypeOf(key).Elem()]; ok {
		return f2(ctx)
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

func RegisterBindingResolver[I any](ctx context.Context, resolver ResolverFunc[I]) context.Context {
	// check if we have a dynamic binding resolver available
	dy := getDynamicBindingResolver(ctx)
	if dy == nil {
		dy = RawBindingResolver{}
	}

	elm := reflect.TypeOf((*I)(nil)).Elem()

	dy[elm] = func(ctx context.Context) (any, error) {
		return resolver(ctx)
	}

	ctx = setDynamicBindingResolver(ctx, dy)

	return ctx
}
