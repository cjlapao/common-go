package serviceprovider

import "fmt"

var globalContainer *ServiceProvider

type Service interface {
	New() interface{}
	Get() interface{}
}

type ServiceLifeCycle int8

const (
	ScopeLifeCycle ServiceLifeCycle = iota
	TransientLifeCycle
)

type ServiceProviderItem struct {
	Service   Service
	LifeCyCle ServiceLifeCycle
}

type ServiceProvider struct {
	services []ServiceProviderItem
}

func RegisterScoped[T Service]() error {
	return Register[T](ScopeLifeCycle)
}

func RegisterTransient[T Service]() error {
	return Register[T](TransientLifeCycle)
}

func Register[T Service](lifeCycle ServiceLifeCycle) error {
	if globalContainer == nil {
		globalContainer = &ServiceProvider{
			services: make([]ServiceProviderItem, 0),
		}
	}

	found := false
	t := *new(T)
	svcType := fmt.Sprintf("%T", t)
	for _, svc := range globalContainer.services {
		xType := fmt.Sprintf("%T", svc.Service)
		if xType == svcType {
			found = true
			break
		}
	}

	if !found {
		service := ServiceProviderItem{
			Service:   t,
			LifeCyCle: lifeCycle,
		}

		globalContainer.services = append(globalContainer.services, service)
	}

	return nil
}

func Get[T Service]() *T {
	t := *new(T)
	svcType := fmt.Sprintf("%T", t)
	for _, svc := range globalContainer.services {
		xType := fmt.Sprintf("%T", svc.Service)
		if xType == svcType {
			switch svc.LifeCyCle {
			case ScopeLifeCycle:
				i := svc.Service.Get()
				return i.(*T)
			case TransientLifeCycle:
				i := svc.Service.New()
				return i.(*T)
			default:
				i := svc.Service.Get()
				return i.(*T)
			}
		}
	}

	return &t
}
