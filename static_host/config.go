package static_host

import (
	"encoding/json"
	"errors"
	"regexp"
)

type config struct {
	Manifest map[string]bool `json:"manifest"`
	Routes   []route         `json:"rules"`
}
type header struct {
	Overwrite *bool  `json:"overwrite,omitempty"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type route struct {
	UseFilesystem *bool          `json:"useFilesystem,omitempty"`
	Source        *regexp.Regexp `json:"source,omitempty"`
	Destination   *string        `json:"destination,omitempty"`
	StatusCode    *uint          `json:"statusCode,omitempty"`
	Headers       *[]header      `json:"headers,omitempty"`
}

type action struct {
	StatusCode  uint
	Headers     []header
	Destination string
}

var (
	errNoMatchedRoute = errors.New("no route matched the path requested")
)

func (c *config) matchRoute(path string) (*action, error) {
	for _, route := range c.Routes {
		if route.UseFilesystem != nil && *route.UseFilesystem{
			if _, ok := c.Manifest[path]; ok {
				return &action{
					StatusCode:  200,
					Headers:     []header{},
					Destination: path,
				}, nil
			}
		}
		if route.Source != nil && route.Source.MatchString(path) {
			newPath := route.Source.ReplaceAllString(path, *route.Destination)
			if _, ok := c.Manifest[newPath]; ok {
				respAction := action{
					StatusCode:  200,
					Destination: newPath,
					Headers:     []header{},
				}
				if route.StatusCode != nil {
					respAction.StatusCode = *route.StatusCode
				}
				if route.Headers != nil {
					respAction.Headers = *route.Headers
				}
				return &respAction, nil
			}
		}
	}
	return nil, errNoMatchedRoute
}

func (c *config) validate() bool {
	if c.Routes != nil {
		isValid := true
		for _, route := range c.Routes {
			isValid = route.validate()
			if !isValid {
				break
			}
		}
		return isValid
	}
	return true
}

func (c *config) MarshalJSON() ([]byte, error) {
	type alias config
	manifestList := make([]string, 0)
	for manifestItem := range c.Manifest {
		manifestList = append(manifestList, manifestItem)
	}
	return json.Marshal(&struct {
		Manifest []string `json:"manifest"`
		*alias
	}{
		Manifest: manifestList,
		alias:    (*alias)(c),
	})
}

func (r *route) MarshalJSON() ([]byte, error) {
	type alias route
	return json.Marshal(&struct {
		Source string `json:"source"`
		*alias
	}{
		Source: r.Source.String(),
		alias: (*alias)(r),
	})
}

func (c *config) UnmarshalJSON(data []byte) error {
	type alias config
	aux := struct {
		Manifest []string `json:"manifest"`
		*alias
	}{
		alias: (*alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	manifestMap := make(map[string]bool)
	for _, manifestItem := range aux.Manifest {
		manifestMap[manifestItem] = true
	}
	c.Manifest = manifestMap
	return nil
}

func (r *route) UnmarshalJSON(data []byte) error {
	type alias route
	aux := struct {
		Source string `json:"source"`
		*alias
	}{
		alias: (*alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if sourceRegex, err := regexp.Compile(aux.Source); err != nil {
		return err
	} else {
		r.Source = sourceRegex
	}
	return nil
}

func (h *header) validate() bool {
	return h.Key != "" && h.Value != ""
}

func (r *route) validate() bool {
	return (
		(r.UseFilesystem != nil && *r.UseFilesystem) &&
			r.Source == nil &&
			r.Destination == nil &&
			r.StatusCode == nil) ||
		(!(r.UseFilesystem != nil && *r.UseFilesystem) &&
			r.Source != nil &&
			r.Destination != nil &&
			*r.Destination != "" &&
			(r.StatusCode == nil ||
				(r.StatusCode != nil &&
					*r.StatusCode >= 200 &&
					*r.StatusCode <= 499)))
}
