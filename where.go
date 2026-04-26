package qb

import (
	"bytes"
	"strings"
)

type whereClause struct {
	wheres []string
	params []interface{}
}

func (w *whereClause) addWhere(condition string, params ...interface{}) {
	w.wheres = append(w.wheres, condition)
	w.params = append(w.params, params...)
}

func (w *whereClause) writeWhere(buf *bytes.Buffer) {
	if len(w.wheres) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(w.wheres, " AND "))
	}
}
