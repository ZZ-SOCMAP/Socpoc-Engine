package build

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/yanmengfei/spoce/library/decode"
	"github.com/yanmengfei/spoce/library/proto"
	"github.com/yanmengfei/spoce/library/utils"
	"gopkg.in/yaml.v2"
)

type PocEvent struct {
	Name   string               `yaml:"name" json:"name"`
	Set    map[string]string    `yaml:"set" json:"set"`
	Rules  []PocRule            `yaml:"rules" json:"rules"`
	Groups map[string][]PocRule `yaml:"groups" json:"groups"`
}

func NewPocEventWithJsonStr(code string) (*PocEvent, error) {
	var object PocEvent
	if err := json.Unmarshal(utils.StrToBytes(code), &object); err != nil {
		return nil, err
	}
	return &object, nil
}

func NewPocEventWithYamlStr(code string) (*PocEvent, error) {
	var object PocEvent
	if err := yaml.Unmarshal(utils.StrToBytes(code), &object); err != nil {
		return nil, err
	}
	return &object, nil
}

func (p PocEvent) ToJsonStr() (string, error) {
	value, err := json.Marshal(&p)
	if err != nil {
		return "", err
	}
	return utils.BytesToStr(value), nil
}

func (p PocEvent) ToYamlStr() (string, error) {
	value, err := yaml.Marshal(&p)
	if err != nil {
		return "", err
	}
	return utils.BytesToStr(value), nil
}

func (p *PocEvent) DecodeSet(env *cel.Env) map[string]interface{} {
	var set = make(map[string]interface{})
	if len(p.Set) > 0 {
		keys := utils.SortMapKeys(p.Set)
		for i := 0; i < len(keys); i++ {
			switch p.Set[keys[i]] {
			case "":
				continue
			case "newReverse()":
				set[keys[i]] = proto.NewReverse()
			default:
				out, err := decode.Evaluate(env, p.Set[keys[i]], set)
				if err == nil {
					switch v := out.Value().(type) {
					case int64:
						set[keys[i]] = int(v)
					default:
						set[keys[i]] = fmt.Sprintf("%v", out)
					}
				}
			}
		}
	}
	return set
}

func (p *PocEvent) DeepcopyRules(source []PocRule) []PocRule {
	rules := make([]PocRule, len(source))
	for i := 0; i < len(rules); i++ {
		rules[i] = source[i]
		rules[i].Headers = make(map[string]string, len(source[i].Headers))
		for key := range source[i].Headers {
			rules[i].Headers[key] = source[i].Headers[key]
		}
	}
	return rules
}

func (p *PocEvent) DeepcopyGroups() map[string][]PocRule {
	var groups = make(map[string][]PocRule, len(p.Groups))
	for key := range p.Groups {
		groups[key] = p.DeepcopyRules(p.Groups[key])
	}
	return groups
}

// PocRule 规则
type PocRule struct {
	Method          string            `yaml:"method" json:"method"`
	Path            string            `yaml:"path" json:"path"`
	Headers         map[string]string `yaml:"headers" json:"headers"`
	Body            string            `yaml:"body" json:"body"`
	FollowRedirects bool              `yaml:"follow_redirects" json:"follow_redirects"`
	Search          string            `yaml:"search" json:"search"`
	Expression      string            `yaml:"expression" json:"expression"`
}

// DecodeSet parse var to set
func (r *PocRule) DecodeSet(set map[string]interface{}) {
	for key, valueSet := range set {
		if _, ok := valueSet.(map[string]string); ok {
			continue
		}
		value := fmt.Sprintf("%v", valueSet)
		for headerKey, headerValue := range r.Headers {
			r.Headers[headerKey] = strings.ReplaceAll(headerValue, "{{"+key+"}}", value)
		}
		r.Path = strings.ReplaceAll(strings.TrimSpace(r.Path), "{{"+key+"}}", value)
		r.Body = strings.ReplaceAll(strings.TrimSpace(r.Body), "{{"+key+"}}", value)
	}
	r.Path = strings.ReplaceAll(r.Path, " ", "%20")
	r.Path = strings.ReplaceAll(r.Path, "+", "%20")
}
