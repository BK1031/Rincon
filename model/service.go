package model

type Service struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Endpoint    string `json:"endpoint" gorm:"unique"`
	HealthCheck string `json:"health_check"`
	UpdatedAt   string `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt   string `json:"created_at" gorm:"autoCreateTime"`
}

func (Service) TableName() string {
	return "service"
}

type ServiceDependency struct {
	ID        string `json:"id" gorm:"primaryKey"`
	ParentID  string `json:"parent_id"`
	ChildID   string `json:"child_id"`
	UpdatedAt string `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt string `json:"created_at" gorm:"autoCreateTime"`
}

func (ServiceDependency) TableName() string {
	return "service_dependency"
}
