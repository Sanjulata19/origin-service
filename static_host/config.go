package static_host

type config struct {
	Rules *[]rule `json:"rules,omitempty"`
}
type header struct {
	key   string
	value string
}

type route struct {
	UseFilesystem *bool   `json:"useFilesystem,omitempty"`
	Source        *string `json:"source,omitempty"`
	Destination   *string `json:"destination,omitempty"`
	StatusCode    *uint   `json:"statusCode,omitempty"`
}

type redirect struct {
	Source      *string `json:"source,omitempty"`
	Destination *string `json:"destination,omitempty"`
	StatusCode  *uint   `json:"statusCode,omitempty"`
}

type rewrite struct {
	Source      *string `json:"source,omitempty"`
	Destination *string `json:"destination,omitempty"`
}

type rule struct {
	Route      *route      `json:"route,omitempty"`
	Redirect   *redirect   `json:"redirect,omitempty"`
	Rewrite    *rewrite    `json:"rewrite,omitempty"`
	SetHeaders *setHeaders `json:"setHeaders,omitempty"`
}

type setHeaders struct {
	Source   *string   `json:"source,omitempty"`
	Headers  *[]header `json:"headers,omitempty"`
	Continue *bool     `json:"continue,omitempty"`
}

func (c *config) validate() bool {
	if c.Rules != nil {
		isValid := true
		for _, rule := range *c.Rules {
			isValid = rule.validate()
			if !isValid {
				break
			}
		}
		return isValid
	}
	return true
}

func (h *header) validate() bool {
	return h.key != "" && h.value != ""
}

func (r *redirect) validate() bool {
	return r.Source != nil &&
		r.Destination != nil &&
		*r.Source != "" &&
		*r.Destination != "" &&
		(r.StatusCode == nil ||
			(r.StatusCode != nil &&
				*r.StatusCode >= 300 &&
				*r.StatusCode <= 399))
}

func (r *rewrite) validate() bool {
	return r.Source != nil &&
		r.Destination != nil &&
		*r.Source != "" &&
		*r.Destination != ""
}

func (r *route) validate() bool {
	return (
		(r.UseFilesystem != nil && *r.UseFilesystem) &&
			r.Source == nil &&
			r.Destination == nil &&
			r.StatusCode == nil) ||
		(!(r.UseFilesystem != nil && *r.UseFilesystem) &&
			r.Source != nil &&
			*r.Source != "" &&
			r.Destination != nil &&
			*r.Destination != "" &&
			(r.StatusCode == nil ||
				(r.StatusCode != nil &&
					*r.StatusCode >= 200 &&
					*r.StatusCode <= 499)))
}

func (r *rule) validate() bool {
	return (r.Route != nil && r.Redirect == nil && r.Rewrite == nil && r.SetHeaders == nil && r.Route.validate()) ||
		(r.Route == nil && r.Redirect != nil && r.Rewrite == nil && r.SetHeaders == nil && r.Redirect.validate()) ||
		(r.Route == nil && r.Redirect == nil && r.Rewrite != nil && r.SetHeaders == nil && r.Rewrite.validate()) ||
		(r.Route == nil && r.Redirect == nil && r.Rewrite == nil && r.SetHeaders != nil && r.SetHeaders.validate())
}

func (s *setHeaders) validate() bool {
	if s.Source != nil && s.Headers != nil && *s.Source != "" {
		isValid := true
		for _, header := range *s.Headers {
			isValid = header.validate()
			if !isValid {
				break
			}
		}
		return isValid
	}
	return false
}
