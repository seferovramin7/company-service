package company

import (
	"company-service/internal/auth"
	"company-service/internal/kafka"
	"company-service/proto"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// CompanyServiceImpl implements the CompanyServiceServer interface.
type CompanyServiceImpl struct {
	proto.UnimplementedCompanyServiceServer
	AuthService   *auth.AuthService
	DB            *sql.DB
	KafkaProducer kafka.Producer // Use the Producer interface for Kafka dependency injection
}

// NewCompanyServiceImpl creates a new CompanyServiceImpl instance with AuthService, DB, and KafkaProducer.
func NewCompanyServiceImpl(authService *auth.AuthService, db *sql.DB, kafkaProducer kafka.Producer) *CompanyServiceImpl {
	return &CompanyServiceImpl{
		AuthService:   authService,
		DB:            db,
		KafkaProducer: kafkaProducer,
	}
}

// CompanyEvent represents the structure of an event message.
type CompanyEvent struct {
	EventType string         `json:"event_type"`
	Company   *proto.Company `json:"company"`
}

// publishEvent sends a company event to Kafka.
func (s *CompanyServiceImpl) publishEvent(ctx context.Context, eventType string, company *proto.Company) {
	event := CompanyEvent{EventType: eventType, Company: company}
	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	companyIDStr := fmt.Sprintf("%d", company.Id)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := s.KafkaProducer.Publish(ctx, companyIDStr, string(eventData)); err != nil {
		log.Printf("Failed to publish %s event for company ID %s: %v", eventType, companyIDStr, err)
	} else {
		log.Printf("Successfully published %s event for company ID %s", eventType, companyIDStr)
	}
}

// CreateCompany handles the creation of a new company.
func (s *CompanyServiceImpl) CreateCompany(ctx context.Context, req *proto.CreateCompanyRequest) (*proto.CompanyID, error) {
	company := req.Company
	query := `
		INSERT INTO companies (name, description, employees, registered, type)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var id int64
	err := s.DB.QueryRowContext(ctx, query, company.Name, company.Description, company.Employees, company.Registered, company.Type).Scan(&id)
	if err != nil {
		log.Printf("Failed to create company: %v", err)
		return nil, err
	}

	company.Id = id
	s.publishEvent(ctx, "CREATE", company)

	return &proto.CompanyID{Id: id}, nil
}

// UpdateCompany handles updating an existing company.
func (s *CompanyServiceImpl) UpdateCompany(ctx context.Context, req *proto.UpdateCompanyRequest) (*proto.CompanyID, error) {
	company := req.Company
	query := `
		UPDATE companies
		SET name = COALESCE(NULLIF($1, ''), name),
		    description = COALESCE(NULLIF($2, ''), description),
		    employees = COALESCE(NULLIF($3::int, 0), employees),
		    registered = COALESCE(NULLIF($4::boolean, FALSE), registered),
		    type = COALESCE(NULLIF($5, ''), type)
		WHERE id = $6
	`
	_, err := s.DB.ExecContext(ctx, query, company.Name, company.Description, company.Employees, company.Registered, company.Type, company.Id)
	if err != nil {
		log.Printf("Failed to update company with id %d: %v", company.Id, err)
		return nil, err
	}

	s.publishEvent(ctx, "UPDATE", company)

	return &proto.CompanyID{Id: company.Id}, nil
}

// DeleteCompany handles deleting a company.
func (s *CompanyServiceImpl) DeleteCompany(ctx context.Context, req *proto.DeleteCompanyRequest) (*proto.CompanyID, error) {
	query := "DELETE FROM companies WHERE id = $1"
	_, err := s.DB.ExecContext(ctx, query, req.Id)
	if err != nil {
		log.Printf("Failed to delete company with id %d: %v", req.Id, err)
		return nil, err
	}

	s.publishEvent(ctx, "DELETE", &proto.Company{Id: req.Id})

	return &proto.CompanyID{Id: req.Id}, nil
}

// GetCompany retrieves a company by its ID.
func (s *CompanyServiceImpl) GetCompany(ctx context.Context, req *proto.CompanyID) (*proto.GetCompanyResponse, error) {
	var company proto.Company
	query := "SELECT id, name, description, employees, registered, type FROM companies WHERE id = $1"
	err := s.DB.QueryRowContext(ctx, query, req.Id).Scan(
		&company.Id,
		&company.Name,
		&company.Description,
		&company.Employees,
		&company.Registered,
		&company.Type,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Company with id %d not found", req.Id)
			return nil, sql.ErrNoRows
		}
		log.Printf("Failed to retrieve company with id %d: %v", req.Id, err)
		return nil, err
	}

	return &proto.GetCompanyResponse{Company: &company}, nil
}

// Login generates a JWT for a given user_id and returns it.
func (s *CompanyServiceImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.AuthService.GenerateToken(req.UserId)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, err
	}

	return &proto.LoginResponse{Token: token}, nil
}
