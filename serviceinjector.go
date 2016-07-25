package river

import "reflect"

type serviceInjector map[reflect.Type]interface{}

// Register registers a new service. Services are identifiable by their types.
// Multiple services of same type should be grouped into a struct,
// and the struct should be registered instead.
func (s *serviceInjector) Register(service interface{}) {
	// don't allow context override
	if reflect.TypeOf(service) == reflect.TypeOf(&Context{}) {
		return
	}
	s.register(service)
}

func (s *serviceInjector) register(service interface{}) {
	if *s == nil {
		*s = make(serviceInjector)
	}
	(*s)[reflect.TypeOf(service)] = service
}

// invoke invokes function f. f must be Func type.
func (s serviceInjector) invoke(f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		// log and return to prevent panic.
		log.println("Cannot invoke non function type")
		return
	}

	args := make([]reflect.Value, reflect.TypeOf(f).NumIn())
	for i := range args {
		argType := reflect.TypeOf(f).In(i)
		if service, ok := s[argType]; ok {
			args[i] = reflect.ValueOf(service)
		} else {
			// set zero value
			args[i] = reflect.Zero(argType)
		}
	}
	reflect.ValueOf(f).Call(args)
}

func copyInjectors(injectors ...serviceInjector) serviceInjector {
	var s = serviceInjector{}
	for i := range injectors {
		for j := range injectors[i] {
			s[j] = injectors[i][j]
		}
	}
	return s
}
