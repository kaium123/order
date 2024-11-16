package sqlxdb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhereBuilder(t *testing.T) {
	sut := NewWhereBuilder("where id>$1", 5, 8.0)

	err := sut.AddCondition("and name like $", "n1")
	require.NoError(t, err)

	now := time.Now()
	err = sut.AddCondition("and dt between $ and $", "t1", now)
	require.NoError(t, err)

	err = sut.AddCondition("and v < $", "t1", "t2")
	require.Error(t, err)

	err = sut.AddCondition("and 1=1")
	require.NoError(t, err)

	gotWhere, gotArgs := sut.Get()

	assert.Equal(t, "where id>$1\nand name like $3\nand dt between $4 and $5\nand 1=1\n", gotWhere)
	assert.EqualValues(t, []interface{}{5, 8.0, "n1", "t1", now}, gotArgs)
}
