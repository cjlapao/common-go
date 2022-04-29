package apiclient

type ApiClientOptions struct {
	Protocol        string
	Host            string
	Path            string
	QueryParameters map[string]string
	Body            interface{}
	Headers         map[string]string
}
