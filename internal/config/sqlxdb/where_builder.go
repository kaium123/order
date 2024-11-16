package sqlxdb

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

var ErrConditionAndArgsCountMismatch = errors.New(
	"condition bind and args count mismatch")

type WhereBuilder struct {
	conditions *bytes.Buffer
	args       []interface{}
}

func NewWhereBuilder(start string, args ...interface{}) *WhereBuilder {
	builder := &WhereBuilder{
		conditions: bytes.NewBuffer([]byte(start)),
		args:       make([]interface{}, len(args)),
	}
	copy(builder.args, args)
	builder.conditions.WriteString("\n")

	return builder
}

// AddCondition adds condition to sql builder. Replaces all $ in condition
// to $1, $2 etc. Checks args count
func (b *WhereBuilder) AddCondition(condition string,
	args ...interface{}) error {

	if len(args) == 0 {
		if strings.Contains(condition, "$") {
			return ErrConditionAndArgsCountMismatch
		}

		b.conditions.WriteString(condition)
		b.conditions.WriteString("\n")
		return nil
	}

	var (
		lastPos  int
		argIndex = len(b.args) + 1
	)

	condBuffer := bytes.NewBuffer(nil)
	for pos, char := range condition {
		if char == '$' {
			condBuffer.WriteString(condition[lastPos : pos+1])
			condBuffer.WriteString(strconv.Itoa(argIndex))
			lastPos = pos + 1
			argIndex++
		}
	}

	if argIndex-len(b.args)-1 != len(args) {
		return ErrConditionAndArgsCountMismatch
	}

	b.conditions.Write(condBuffer.Bytes())
	b.conditions.WriteString("\n")
	b.args = append(b.args, args...)

	return nil
}

// MustAddCondition adds condition to sql builder. Panics on error.
func (b *WhereBuilder) MustAddCondition(condition string, args ...interface{}) {
	err := b.AddCondition(condition, args...)
	if err != nil {
		panic(err)
	}
}

// Get returns current where and args
func (b *WhereBuilder) Get() (where string, args []interface{}) {
	return b.conditions.String(), b.args
}
