package commands

type ExecuteResponse struct {
	StdOut    string
	StdErr    string
	ErrorCode string
}

func (r ExecuteResponse) GetAllOutputs() string {
	result := ""
	if r.StdOut != "" {
		result += r.StdOut
	}
	if r.StdErr != "" {
		result += r.StdErr
	}

	return result
}
