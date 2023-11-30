package snake

import (
	"context"

	"github.com/go-faster/errors"
)

func ResolveAllShared(ctx context.Context, names []string, fmap FMap, binder *Binder) (*Binder, error) {

	for _, v := range names {
		var err error
		var resolver Resolver
		if resolver = fmap(v); resolver == nil {
			return nil, errors.Errorf("missing resolver for %q", v)
		}

		if MenthodIsShared(resolver) {
			binder, err = findArgumentsRaw(v, fmap, binder)
			if err != nil {
				return nil, err
			}
		}

	}
	return binder, nil
}

func WrapWithMiddleware(base MiddlewareFunc, middlewares ...Middleware) MiddlewareFunc {

	for _, v := range middlewares {
		base = v.Wrap(base)
	}
	return base
}

func RunResolvingArguments(outputHandler OutputHandler, fmap FMap, str string, binder *Binder, middlewares ...Middleware) error {
	// always resolve context.Context first
	_, err := findArgumentsRaw("context.Context", fmap, binder)
	if err != nil {
		return err
	}

	base := func(ctx context.Context) error {
		defer func() {
			delete(binder.bindings, str)
		}()

		binder, err := findArgumentsRaw(str, fmap, binder)
		if err != nil {
			return err
		}

		out := binder.bindings[str]

		if out == nil {
			return errors.Errorf("missing resolver for %q", str)
		}

		result := binder.bindings[str].Interface()

		if out, ok := result.(Output); ok {
			return HandleOutput(ctx, outputHandler, out)
		}

		return nil
	}

	wrp := WrapWithMiddleware(base, middlewares...)

	ctx := binder.bindings["context.Context"].Interface().(context.Context)

	err = wrp(ctx)
	if err != nil {
		return err
	}

	return nil
}