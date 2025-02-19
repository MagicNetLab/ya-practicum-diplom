package local

type Repository struct {
	baseFolder string
}

func (r Repository) Close() error {
	return nil
}
