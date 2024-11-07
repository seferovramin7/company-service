package company

import (
	"company-service/internal/auth"
	"company-service/internal/kafka"
	"company-service/proto"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCompany(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	kafkaProducer := &kafka.KafkaProducerMock{}
	authService := auth.NewAuthService("test-secret")

	// Expecting an INSERT statement
	mock.ExpectQuery("INSERT INTO companies").
		WithArgs("Test Co", "A sample company", 50, true, "Corporation").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	service := NewCompanyServiceImpl(authService, db, kafkaProducer)

	req := &proto.CreateCompanyRequest{
		Company: &proto.Company{
			Name:        "Test Co",
			Description: "A sample company",
			Employees:   50,
			Registered:  true,
			Type:        "Corporation",
		},
	}

	resp, err := service.CreateCompany(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp.Company)
	assert.Equal(t, int64(1), resp.Company.Id)
	assert.Len(t, kafkaProducer.PublishedMessages, 1)
	assert.Contains(t, kafkaProducer.PublishedMessages[0], `"event_type":"CREATE"`)
}

func TestUpdateCompany(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	kafkaProducer := &kafka.KafkaProducerMock{}
	authService := auth.NewAuthService("test-secret")

	// Expecting an UPDATE statement
	mock.ExpectExec("UPDATE companies").
		WithArgs("Updated Co", "Updated description", 100, false, "LLC", int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	service := NewCompanyServiceImpl(authService, db, kafkaProducer)

	req := &proto.UpdateCompanyRequest{
		Company: &proto.Company{
			Id:          1,
			Name:        "Updated Co",
			Description: "Updated description",
			Employees:   100,
			Registered:  false,
			Type:        "LLC",
		},
	}

	resp, err := service.UpdateCompany(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp.Company)
	assert.Equal(t, int64(1), resp.Company.Id)
	assert.Len(t, kafkaProducer.PublishedMessages, 1)
	assert.Contains(t, kafkaProducer.PublishedMessages[0], `"event_type":"UPDATE"`)
}

func TestDeleteCompany(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	kafkaProducer := &kafka.KafkaProducerMock{}
	authService := auth.NewAuthService("test-secret")

	// Expecting a DELETE statement
	mock.ExpectExec("DELETE FROM companies WHERE id = \\$1").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	service := NewCompanyServiceImpl(authService, db, kafkaProducer)

	req := &proto.DeleteCompanyRequest{Id: 1}

	resp, err := service.DeleteCompany(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Id)
	assert.Len(t, kafkaProducer.PublishedMessages, 1)
	assert.Contains(t, kafkaProducer.PublishedMessages[0], `"event_type":"DELETE"`)
}

func TestGetCompany(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	kafkaProducer := &kafka.KafkaProducerMock{}
	authService := auth.NewAuthService("test-secret")

	// Expecting a SELECT statement
	mock.ExpectQuery("SELECT id, name, description, employees, registered, type FROM companies WHERE id = \\$1").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "employees", "registered", "type"}).
			AddRow(1, "Test Co", "A sample company", 50, true, "Corporation"))

	service := NewCompanyServiceImpl(authService, db, kafkaProducer)

	req := &proto.CompanyID{Id: 1}

	resp, err := service.GetCompany(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp.Company)
	assert.Equal(t, req.Id, resp.Company.Id)
	assert.Equal(t, "Test Co", resp.Company.Name)
	assert.Equal(t, "A sample company", resp.Company.Description)
	assert.Equal(t, int32(50), resp.Company.Employees)
	assert.True(t, resp.Company.Registered)
	assert.Equal(t, "Corporation", resp.Company.Type)
}

func TestLogin(t *testing.T) {
	authService := auth.NewAuthService("test-secret")
	service := NewCompanyServiceImpl(authService, nil, nil)

	req := &proto.LoginRequest{UserId: 1}

	resp, err := service.Login(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}
