package company

import "company-service/proto"

type Company struct {
	ID          int64
	Name        string
	Description string
	Employees   int
	Registered  bool
	Type        string
}

// ToProto converts Company struct to protobuf Company message
func (c *Company) ToProto() *proto.Company {
	return &proto.Company{
		Id:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Employees:   int32(c.Employees),
		Registered:  c.Registered,
		Type:        c.Type,
	}
}

// FromProto converts protobuf Company message to Company struct
func FromProto(protoCompany *proto.Company) *Company {
	return &Company{
		ID:          protoCompany.Id,
		Name:        protoCompany.Name,
		Description: protoCompany.Description,
		Employees:   int(protoCompany.Employees),
		Registered:  protoCompany.Registered,
		Type:        protoCompany.Type,
	}
}
