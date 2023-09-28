package snake

import (
	"context"
	"reflect"
)

func ResolveBindingsFromProvider(ctx context.Context, rf reflect.Value) (context.Context, error) {

	for _, pt := range listOfArgs(rf.Type()) {
		// if pt.Kind() == reflect.Ptr {
		// 	pt = pt.Elem()
		// }

		rslv, ok := getResolvers(ctx)[pt]
		if !ok {
			continue
		}

		p, err := rslv(ctx)
		if err != nil {
			return ctx, err
		}

		if reflect.TypeOf(p.Interface()).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			// if the provider returns a context - meaning the dyanmic context binding resolver was set
			// we need to merge any bindings that might have been set
			// ctx = mergeResolversKeepingFirst(ctx, p.(context.Context))

			crb := p.Interface().(context.Context)

			// this is the context resolver binding, we need to process it
			ctx = mergeBindingKeepingFirst(ctx, crb)

		}
		ctx = bind(ctx, pt, p)
		break

	}

	return ctx, nil
}

func setResolvers(ctx context.Context, provider resolvers) context.Context {
	return context.WithValue(ctx, &resolverKeyT{}, provider)
}

func getResolvers(ctx context.Context) resolvers {
	p, ok := ctx.Value(&resolverKeyT{}).(resolvers)
	if ok {
		return p
	}
	return resolvers{}
}

func setFlagBindings(ctx context.Context, provider flagbindings) context.Context {
	return context.WithValue(ctx, &flagbindingsKeyT{}, provider)
}

func getFlagBindings(ctx context.Context) flagbindings {
	p, ok := ctx.Value(&flagbindingsKeyT{}).(flagbindings)
	if ok {
		return p
	}
	return flagbindings{}
}

func RegisterBindingResolver[I any](ctx context.Context, res typedResolver[I], f ...flagbinding) context.Context {
	// check if we have a dynamic binding resolver available
	dy := getResolvers(ctx)

	elm := reflect.TypeOf((*I)(nil)).Elem()

	dy[elm] = res.asResolver()

	ctx = setResolvers(ctx, dy)

	for _, fbb := range f {
		fb := getFlagBindings(ctx)
		fb[elm] = fbb
		ctx = setFlagBindings(ctx, fb)
	}

	return ctx
}
