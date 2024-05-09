package Register_go

type ServiceProvider struct {
	Registry
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (sp ServiceProvider) FindService(name string) []ServiceInfo {
	_, ok := GlobalRegistry.nameToId[name]
	if !ok {
		return nil
	} else {
		idList := GlobalRegistry.nameToId[name]
		serviceList := make([]ServiceInfo, len(idList), len(idList))
		for i, id := range idList {
			serviceList[i] = GlobalRegistry.services[id]
		}
		return serviceList
	}
}
