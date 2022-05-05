package apiclient

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/guard"
)

type ApiClientOptions struct {
	Method          ApiClientMethod
	Authorization   ApiClientAuthorization
	Protocol        string
	Port            int
	Host            string
	Path            string
	QueryParameters map[string]string
	Body            ApiClientBody
	Headers         map[string]string
}

func (options ApiClientOptions) Url() (*url.URL, error) {
	if err := guard.EmptyOrNil(options.Host); err != nil {
		return nil, err
	}

	// Validating if host does not have the full url
	if strings.HasPrefix(options.Host, "http://") {
		options.Protocol = "http"
		options.Host = strings.TrimLeft(options.Host, "http://")
	}
	if strings.HasPrefix(options.Host, "https://") {
		options.Protocol = "https"
		options.Host = strings.TrimLeft(options.Host, "https://")
	}

	if strings.ContainsAny(options.Host, ":") {
		hostParts := strings.Split(options.Host, ":")
		if len(hostParts) > 1 {
			portNumber, portErr := strconv.Atoi(hostParts[1])
			if portErr == nil {
				options.Port = portNumber
			}
		}
		options.Host = hostParts[0]
	}

	if err := guard.EmptyOrNil(options.Protocol); err != nil {
		return nil, err
	}

	if strings.HasSuffix(options.Host, "/") {
		options.Host = strings.TrimRight(options.Host, "/")
	}

	rawUrl := fmt.Sprintf("%v://%v", options.Protocol, options.Host)

	if options.Port > 0 {
		rawUrl = fmt.Sprintf("%v:%v", rawUrl, options.Port)
	}

	if options.Path != "" {
		if !strings.HasPrefix(options.Path, "/") {
			options.Path = fmt.Sprintf("/%v", options.Path)
		}
		rawUrl = fmt.Sprintf("%v%v", rawUrl, options.Path)
		rawUrl = strings.TrimRight(rawUrl, "/")
	}

	if options.QueryParameters != nil && len(options.QueryParameters) > 0 {
		count := 0
		for key, value := range options.QueryParameters {
			if count == 0 {
				rawUrl = fmt.Sprintf("%v?", rawUrl)
			}
			rawUrl = fmt.Sprintf("%v%v=%v", rawUrl, key, value)
			count += 1
			if count < len(options.QueryParameters) {
				rawUrl = fmt.Sprintf("%v&", rawUrl)
			}
		}
	}

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	return url, nil
}
