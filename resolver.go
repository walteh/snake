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

func RunResolvingArguments(outputHandler OutputHandler, fmap FMap, str string, binder *Binder) error {
	// always resolve context.Context first
	binder, err := findArgumentsRaw("context.Context", fmap, binder)
	if err != nil {
		return err
	}

	binder, err = findArgumentsRaw(str, fmap, binder)
	if err != nil {
		return err
	}

	out := binder.bindings[str]
	ctx := binder.bindings["context.Context"].Interface().(context.Context)

	if out == nil {
		return errors.Errorf("missing resolver for %q", str)
	}

	result := binder.bindings[str].Interface()

	if out, ok := result.(Output); ok {
		return HandleOutput(ctx, outputHandler, out)
	}

	return nil
}
