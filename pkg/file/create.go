package file

import "context"

func Create(ctx context.Context, r repo, f *File) error {
	err := r.Insert(ctx, f)
	if err != nil {
		return err
	}

	return nil
}
