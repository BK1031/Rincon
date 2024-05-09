package service

import (
	"rincon/database"
	"rincon/model"
	"rincon/utils"
	"time"
)

func GetAllServices() []model.Service {
	var services []model.Service
	services = database.Local.Services
	return services
}

func GetServiceByID(id int) model.Service {
	var service model.Service
	for _, s := range database.Local.Services {
		if s.ID == id {
			service = s
			break
		}
	}
	return service
}

func GetServicesByName(name string) []model.Service {
	var services []model.Service
	for _, s := range database.Local.Services {
		if s.Name == name {
			services = append(services, s)
		}
	}
	return services
}

func GetServiceByEndpoint(endpoint string) model.Service {
	var service model.Service
	for _, s := range database.Local.Services {
		if s.Endpoint == endpoint {
			service = s
			break
		}
	}
	return service
}

func CreateService(service model.Service) error {
	existing := GetServiceByEndpoint(service.Endpoint)
	if existing.Endpoint != "" {
		existing.UpdatedAt = time.Now()
		for i, s := range database.Local.Services {
			if s.ID == existing.ID {
				database.Local.Services[i] = existing
				break
			}
		}
	} else {
		service.ID = utils.GenerateID(0)
		service.UpdatedAt = time.Now()
		service.CreatedAt = time.Now()
		database.Local.Services = append(database.Local.Services, service)
	}
	return nil
}
