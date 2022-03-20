package mock

type Differ struct {
	DiffOutput []string
	DiffError  error
	DiffInputs []DiffInput
}

func (d *Differ) Diff(from, to string) ([]string, error) {
	if d.DiffError != nil {
		return nil, d.DiffError
	}

	d.DiffInputs = append(d.DiffInputs, DiffInput{from: from, to: to})
	return d.DiffOutput, nil
}

type DiffInput struct {
	from string
	to   string
}
