package proto

import (
	"net/url"
	"strings"
)

// SetupURL URL to UrlType
func SetupURL(u *url.URL) *UrlType {
	return &UrlType{
		Scheme:   u.Scheme,
		Domain:   u.Hostname(),
		Host:     u.Host,
		Port:     u.Port(),
		Path:     u.EscapedPath(),
		Query:    u.RawQuery,
		Fragment: u.Fragment,
	}
}

// ToString UrlType to string
func (x *UrlType) ToString() string {
	var buf strings.Builder
	if x.Scheme != "" {
		buf.WriteString(x.Scheme)
		buf.WriteByte(':')
	}
	if x.Scheme != "" || x.Host != "" {
		if x.Host != "" || x.Path != "" {
			buf.WriteString("//")
		}
		if h := x.Host; h != "" {
			buf.WriteString(x.Host)
		}
	}
	path := x.Path
	if path != "" && path[0] != '/' && x.Host != "" {
		buf.WriteByte('/')
	}
	if buf.Len() == 0 {
		if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
			buf.WriteString("./")
		}
	}
	buf.WriteString(path)
	if x.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(x.Query)
	}
	if x.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(x.Fragment)
	}
	return buf.String()
}
