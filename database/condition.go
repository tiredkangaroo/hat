package database

import (
	"fmt"
	"slices"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	OperatorAND = "AND"
	OperatorOR  = "OR"
	OperatorEQ  = "equals"
	OperatorCT  = "contains"
)

var (
	ErrInvalidCondition = fmt.Errorf("invalid condition")
)

type Condition struct {
	Operator   string      `json:"op"` // "AND", "OR", "equals", "contains", ...
	Field      string      `json:"field,omitempty"`
	Value      any         `json:"value,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"` // sub-conditions for AND/OR
}

type Context struct {
	Device     *Device              // device id of a request (nil if not handling a request)
	RequestCtx *fasthttp.RequestCtx // fasthttp request context (nil if not handling a request)
}

func (ctx *Context) Get(field string) (any, error) {
	if strings.HasPrefix(field, "ctx-") && ctx.RequestCtx == nil {
		return nil, fmt.Errorf("value not available for field: %s", field)
	}
	if strings.HasPrefix(field, "device-") && ctx.Device == nil {
		return nil, fmt.Errorf("value not available for field: %s", field)
	}
	switch field {
	case "device-id":
		return ctx.Device.ID, nil
	case "ctx-host":
		return b2s(ctx.RequestCtx.Host()), nil
	case "ctx-method":
		return b2s(ctx.RequestCtx.Method()), nil
	case "ctx-path":
		return b2s(ctx.RequestCtx.Path()), nil
	case "ctx-body":
		return b2s(ctx.RequestCtx.Request.Body()), nil
	}
	return nil, fmt.Errorf("unknown field: %s", field)
}

func (c *Condition) Evaluate(ctx *Context) (bool, error) {
	switch c.Operator {
	case OperatorAND:
		for _, cond := range c.Conditions {
			ev, err := cond.Evaluate(ctx)
			if !ev || err != nil {
				return false, err
			}
		}
		return true, nil
	case OperatorOR:
		for _, cond := range c.Conditions {
			result, err := cond.Evaluate(ctx)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil
	case OperatorEQ:
		v, err := ctx.Get(c.Field)
		if err != nil {
			return false, err
		}
		return v == c.Value, nil
	case OperatorCT:
		v, err := ctx.Get(c.Field)
		if err != nil {
			return false, err
		}
		return handleContains(v, c.Value)
	}
	return false, fmt.Errorf("unknown operator: %s", c.Operator)
}

func handleContains(x, subx any) (bool, error) {
	// assume the types are the same
	if x == subx {
		return true, nil
	}
	switch x := x.(type) {
	case string:
		return strings.Contains(x, subx.(string)), nil
	case []string:
		return slices.Contains(x, subx.(string)), nil
	case map[string]any:
		_, ok := x[subx.(string)]
		return ok, nil
	}
	return false, fmt.Errorf("cannot handle contains for type: %T", x)
}
