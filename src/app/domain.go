// Package app represents app's domain layer
// the domain layer is business-specific and application-agnostic
package app

import (
	"net/http"
	"time"
)

// ErrorHandler interface
type ErrorHandler interface {
	Handle(http.ResponseWriter, error)
}

// MailSender interface
type MailSender interface {
	Send(to []string, subject string, body []byte) error
}

// CDNUploader interface
type CDNUploader interface {
	Upload(string) (*http.Response, error)
}

// Model is model trait
type Model struct {
	ID        int       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ModelSoftDelete model trait that has DeletedAt field
type ModelSoftDelete struct {
	Model
	DeletedAt *time.Time `json:"deletedAt"`
}

type Account struct {
	Model
	Name string `json:"name"`
	Desc string `json:"desc"`
	Type string `json:"type"`
}

type TransactionTag struct {
	Model
	Name string `json:"name"`
}

type Transaction struct {
	Model
	UserID    int              `json:"-"`
	AccountID int              `json:"-"`
	Account   Account          `json:"account"`
	Amount    float32          `json:"amount"`
	Desc      string           `json:"desc"`
	Tags      []TransactionTag `json:"tags,omitempty" gorm:"many2many:pivot_transaction_tag;"`
	Date      time.Time        `json:"date"`
}

func (t *Transaction) AddTag(cs ...string) {
	for _, c := range cs {
		t.addTag(c)
	}
}

func (t *Transaction) addTag(tag string) {
	for _, v := range t.Tags {
		if v.Name == tag {
			return
		}
	}
	t.Tags = append(t.Tags, TransactionTag{Name: tag})
}

type Image struct {
	Model
	PublicID     string `json:"publicId" gorm:"unique_index"`
	ResourceType string `json:"resourceType"`
}
