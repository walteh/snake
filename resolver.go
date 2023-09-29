package snake

import (
	"context"
	"reflect"
)

func ResolveBindingsFromProvider(ctx context.Context, rf reflect.Value) (context.Context, error) {

	loa := listOfArgs(rf.Type())

	if len(loa) > 1 {
		for i, pt := range loa {
			if i == 0 {
				continue
			}
			if pt.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
				// move the context to the first position
				// this is so that subsequent bindings can be resolved with the updated context
				loa[0], loa[i] = loa[i], loa[0]
				break
			}
		}
	}

	for _, pt := range loa {
		rslv, ok := getResolvers(ctx)[pt]
		if !ok {
			if pt.Kind() == reflect.Ptr {
				pt = pt.Elem()
			}
			// check if we have a flag binding for this type
			rslv, ok = getResolvers(ctx)[pt]
			if !ok {
				continue
			}
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
			// we favor the context returned from the resolver, as it also might have been modified
			// for example, if the resolver returns a context with a zerolog logger, we want to keep that
			ctx = mergeBindingKeepingFirst(crb, ctx)

		}
		ctx = bind(ctx, pt, p)

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
