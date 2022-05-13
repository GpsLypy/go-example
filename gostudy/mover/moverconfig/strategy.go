package moverconfig

import (
	"app/tools/excavator/bformat"
	"log"
)

const (
	StrategyRowInsert = 1 << iota
	StrategyRowUpdate
	StrategyRowDelete
)

// StrategyConfig contains strategy config
// InterestRow is the interest row events, 0 means all, others are with bit or
// InterestFields is the interest fields, null means all, empty means only contains
// the primary key field, all fields without primary keys will auto contain primary
// key.
type StrategyConfig struct {
	InterestRow    int      `json:"interestRow"`
	InterestFields []string `json:"interestFields"`
	interestFields map[string]struct{}
}

func (sc *StrategyConfig) Init() {
	if nil == sc.InterestFields {
		return
	}
	sc.interestFields = make(map[string]struct{})
	for _, v := range sc.InterestFields {
		sc.interestFields[v] = struct{}{}
	}
}

func (sc *StrategyConfig) clone() *StrategyConfig {
	cl := &StrategyConfig{}
	cl.InterestRow = sc.InterestRow
	cl.InterestFields = make([]string, len(sc.InterestFields))
	copy(cl.InterestFields, sc.InterestFields)
	return cl
}

func (sc *StrategyConfig) Match(row *bformat.RowEvent) bool {
	if 0 == sc.InterestRow {
		return true
	}
	if row.DataType == bformat.DataTypeQuery {
		return true
	}

	var rowFlag int
	if row.DataType == bformat.DataTypeInsert {
		rowFlag |= StrategyRowInsert
	} else if row.DataType == bformat.DataTypeUpdate {
		rowFlag |= StrategyRowUpdate
	} else if row.DataType == bformat.DataTypeDelete {
		rowFlag |= StrategyRowDelete
	}
	return (rowFlag & sc.InterestRow) != 0
}

// DoFilter implements fieldFilter interface
func (sc *StrategyConfig) DoFilter(table *bformat.TableStructEvent,
	row *bformat.RowEvent) (*bformat.TableStructEvent,
	*bformat.RowEvent) {
	if nil == sc.InterestFields {
		return table, row
	}
	if row.DataType == bformat.DataTypeDelete {
		if len(row.Datas) != len(table.Table.Columns) {
			log.Panicf("Mismatch row data count %d and table column count %d",
				len(row.Datas), len(table.Table.Columns))
		}
		// Create clone table
		var newTable bformat.TableStructEvent
		newTable.Table = &bformat.TableDef{}
		newTable.Table.Schema = table.Table.Schema
		newTable.Table.Name = table.Table.Name
		newTable.Table.Columns = make([]*bformat.ColumnDef, 0, len(table.Table.Columns))
		// Create clone row
		var newRow bformat.RowEvent
		newRow.DataType = row.DataType
		newRow.Datas = make([]interface{}, 0, len(row.Datas))
		// Do filter
		for i, col := range table.Table.Columns {
			if _, ok := sc.interestFields[col.Name]; !ok {
				if !col.PrimaryKey {
					continue
				}
			}
			newTable.Table.Columns = append(newTable.Table.Columns, col)
			newRow.Datas = append(newRow.Datas, row.Datas[i])
		}
		return &newTable, &newRow
	}
	// Others are not supported now
	return table, row
}
