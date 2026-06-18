package core

import "strings"

type Rule struct {
	Server string   `json:"server"`
	File   []string `json:"file"`
	Extra  []string `json:"extra"`
	Output []string `json:"output"`
	domain []string
}

func (r *Rule) CanPAC() bool {
	return r.checkOutput(`pac`)
}

func (r *Rule) CanClash() bool {
	return r.checkOutput(`clash`)
}

func (r *Rule) checkOutput(check string) bool {
	if len(r.Output) == 0 {
		return true
	}
	for _, s := range r.Output {
		if strings.ToLower(s) == check {
			return true
		}
	}
	return false
}
