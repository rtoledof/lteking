package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Application struct {
	ID     string          `json:"id" bson:"_id"`
	Name   string          `json:"name" bson:"name"`
	Type   ApplicationType `json:"type" bson:"type"`
	Client string          `json:"client" bson:"client"`
	Secret string          `json:"secret" bson:"secret"`
}

type ApplicationFilter struct {
	Ids    []*string `json:"ids,omitempty"`
	Limit  int       `json:"limit,omitempty"`
	Token  *string   `json:"token,omitempty"`
	Name   *string   `json:"name,omitempty"`
	Client *string   `json:"client,omitempty"`
}

type ApplicationList struct {
	Token *string        `json:"token,omitempty"`
	Data  []*Application `json:"data,omitempty"`
}

type ApplicationRequest struct {
	Name   string          `json:"name"`
	Type   ApplicationType `json:"type"`
	Client string          `json:"client"`
	Secret string          `json:"secret"`
}

type ApplicationService interface {
	FindByClient(ctx context.Context, client string) (*Application, error)
	FindByID(ctx context.Context, input string) (*Application, error)
	FindApplications(ctx context.Context, input *ApplicationFilter) (*ApplicationList, error)
	CreateApplication(ctx context.Context, input ApplicationRequest) (*Application, error)
	UpdateApplicationCredentials(ctx context.Context, application string) (*Application, error)
}

type ApplicationType string

const (
	ApplicationTypeRider  ApplicationType = "RIDER"
	ApplicationTypeDriver ApplicationType = "DRIVER"
)

var AllApplicationType = []ApplicationType{
	ApplicationTypeRider,
	ApplicationTypeDriver,
}

func (e ApplicationType) IsValid() bool {
	switch e {
	case ApplicationTypeRider, ApplicationTypeDriver:
		return true
	}
	return false
}

func (e ApplicationType) String() string {
	return string(e)
}

func (e *ApplicationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ApplicationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ApplicationType", str)
	}
	return nil
}

func (e ApplicationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
