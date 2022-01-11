package icinga

import (
	"encoding/json"
	"fmt"
)

func (s Service) name() string {
	return s.Name
}

func (s Service) path() string {
	return "/objects/services/" + s.Name
}

func (s Service) attrs() map[string]interface{} {
	m := make(map[string]interface{})
	m["display_name"] = s.DisplayName
	return m
}

// Service represents a Service object.
type Service struct {
	Name         string `json:"__name"`
	Groups       []string
	State        ServiceState
	CheckCommand string `json:"check_command"`
	DisplayName  string `json:"display_name:"`
}

type ServiceState int

const (
	ServiceOK ServiceState = 0 + iota
	ServiceWarning
	ServiceCritical
	ServiceUnknown
)

func (s ServiceState) String() string {
	switch s {
	case ServiceOK:
		return "ServiceOK"
	case ServiceWarning:
		return "ServiceWarning"
	case ServiceCritical:
		return "ServiceCritical"
	case ServiceUnknown:
		return "ServiceUnknown"
	}
	return "unhandled service state"
}

func (s Service) MarshalJSON() ([]byte, error) {
	attrs := make(map[string]interface{})
	if len(s.Groups) > 0 {
		attrs["groups"] = s.Groups
	}
	attrs["check_command"] = s.CheckCommand
	attrs["display_name"] = s.DisplayName
	jservice := &struct {
		Attrs map[string]interface{} `json:"attrs"`
	}{Attrs: attrs}
	return json.Marshal(jservice)
}

func (c *Client) CreateService(service Service) error {
	if err := c.createObject(service); err != nil {
		return fmt.Errorf("create service %s: %w", service.Name, err)
	}
	return nil
}

func (c *Client) LookupService(name string) (Service, error) {
	obj, err := c.lookupObject("/objects/services/" + name)
	if err != nil {
		return Service{}, fmt.Errorf("lookup %s: %w", name, err)
	}
	v, ok := obj.(Service)
	if !ok {
		return Service{}, fmt.Errorf("lookup %s: result type %T is not service", name, obj)
	}
	return v, nil
}