//4、复合主键是string型时任务派发逻辑【待解决】
func mysqlDispatchTaskByPriString(dsp moverconfig.DataSourcePair, taskChan *channels.InfiniteChannel, fieldList []string, priKey string, priKeyTyp int) error {
	// get db conn
	dbConn, _, err := getMysqlDB4Ds(dsp.From, "")
	if nil != err {
		return errors.Trace(err)
	}

	var conditions []string
	if "" != dsp.FromWhere {
		conditions = append(conditions, dsp.FromWhere)
	}

	bp := findBreakPoint(getNameByHostAndPort(dsp.From.Endpoints[0].Host, dsp.From.Endpoints[0].Port, dsp.From.DBName), dsp.From.TableName)
	if nil != bp && bp.PrimaryKeyType == priKeyTyp {
		bpCond := fmt.Sprintf("`%s` > '%s'", bp.PrimaryKey, bp.StrEndId)
		conditions = append(conditions, bpCond)
	}

	// get table sum & maxid
	var rowSum, rowCount sql.NullInt64
	var rowMin, rowMax, rowCur sql.NullString
	var xSql, countSql string
	if len(conditions) == 0 {
		xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s`", priKey, priKey, dsp.From.DBName, dsp.From.TableName)
	} else {
		xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s` WHERE %s", priKey, priKey, dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "))
	}

	if len(conditions) == 0 {
		countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, priKey)
	} else {
		countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` WHERE %s group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "), priKey)
	}

	LogDebug("Sql: %s", xSql)

	row := dbConn.QueryRow(xSql)
	err = row.Scan(&rowSum, &rowMin, &rowMax)
	if nil != err {
		logger.Errorf("mysqlDispatchTaskByString-sql=%s-err=%v", xSql, err)
		return errors.Trace(err)
	}

	LogInfo("Sum=%v, min=%v, max=%v, dsp: %v", rowSum, rowMin, rowMax, dsp)
	resetAppendTableStat(TableStat{
		HostPort:  dsp.From.Endpoints[0].Host + ":" + strconv.Itoa(dsp.From.Endpoints[0].Port),
		DbName:    dsp.From.DBName,
		TableName: dsp.From.TableName,
		SumRows:   rowSum.Int64,
		Sql:       xSql,
	})

	if rowSum.Int64 == 0 {
		return nil
	}

	rowC := dbConn.QueryRow(countSql)
	err = rowC.Scan(&rowCount)
	if nil != err {
		logger.Errorf("mysqlDispatchTaskByString-countSqlsql=%s-err=%v", countSql, err)
		return errors.Trace(err)
	}

	var sum int64
	var min, max string
	sum = rowSum.Int64
	if rowMin.Valid {
		min = rowMin.String
	}
	if rowMax.Valid {
		max = rowMax.String
	}

	// split job
	var MyDataSizePerTask int64 = DataSizePerTask
	if dsp.Dest.IsSharding {
		MyDataSizePerTask = DataSizePerTask / 10
		if MyDataSizePerTask <= 0 {
			MyDataSizePerTask = 1
		}
	}

	// for qatest
	{
		if getConfig().RowsPerTask > 0 {
			MyDataSizePerTask = int64(getConfig().RowsPerTask)
		}
	}

	if MyDataSizePerTask < rowCount.Int64+1 {
		MyDataSizePerTask = rowCount.Int64 + 1
	}

	var taskSize int
	if sum%MyDataSizePerTask == 0 {
		taskSize = int(sum / MyDataSizePerTask)
	} else {
		taskSize = int(sum/MyDataSizePerTask + 1)
	}

	if 0 == taskSize {
		var fields string
		if "" == dsp.FromField {
			fields = "`" + strings.Join(fieldList, "`,`") + "`"
		} else {
			fields = priKey + "," + dsp.FromField
		}

		empTask := Task{
			FromDBType:     dsp.From.DBType,
			FromEndpoints:  dsp.From.Endpoints,
			FromEndpoint:   dsp.From.Endpoints[0],
			FromDBName:     dsp.From.DBName,
			FromTable:      dsp.From.TableName,
			FromField:      fields,
			FromWhere:      dsp.FromWhere,
			OrgDBName:      dsp.From.OrgDBName,
			OrgTbName:      dsp.From.OrgTbName,
			PrimaryKey:     priKey,
			PrimaryKeyType: priKeyTyp,
			DestTo:         dsp.DestTo,
			DestMQInf:      dsp.DestMQ,
			ClosedInterval: 0,
			StartId:        0,
			EndId:          0,
		}

		seelog.Debugf("Empty task: %v", empTask)
		taskChan.In() <- empTask
		atomic.AddInt32(&currState.SumTasks, 1)

		return nil
	}

	var task Task
	task.IsSelfPri = true
	task.FromDBType = dsp.From.DBType
	task.FromIsShard = dsp.From.IsSharding
	task.FromEndpoint = dsp.From.Endpoints[0]
	appendStateFromHost(task.FromEndpoint.Host + ":" + strconv.Itoa(task.FromEndpoint.Port))
	task.FromEndpoints = dsp.From.Endpoints
	task.FromDBName = dsp.From.DBName
	task.FromTable = dsp.From.TableName
	task.OrgDBName = dsp.From.OrgDBName
	task.OrgTbName = dsp.From.OrgTbName
	if "" == dsp.FromField {
		task.FromField = "`" + strings.Join(fieldList, "`,`") + "`"
	} else {
		if getConfig().Target == TargetDatabase {
			task.FromField = dsp.FromField
		} else {
			task.FromField = priKey + "," + dsp.FromField
		}
	}
	task.FromWhere = dsp.FromWhere
	task.PrimaryKey = priKey
	task.PrimaryKeyType = priKeyTyp
	task.DestTo = dsp.DestTo

	if dsp.DestTo == 0 {
		task.DestDBType = dsp.Dest.DBType
		task.DestEndpoint = dsp.Dest.Endpoints[0]
		task.DestDBName = dsp.Dest.DBName
		task.DestTable = dsp.Dest.TableName
		task.DestField = fieldList
		task.DestDBIsSharding = dsp.Dest.IsSharding
	} else {
		task.DestMQInf = dsp.DestMQ
	}
	task.RowCountLimit = int(rowCount.Int64) + 1

	destEndpointSize := len(dsp.Dest.Endpoints)
	curId := min

	for i := 0; i < taskSize; i++ {
		if strings.Compare(curId, max) > 0 {
			LogWarn("Premature termination, curId=%s, maxId=%s, dsp: %v", curId, max, dsp)
			break
		}

		task.ClosedInterval = 0
		if 0 == i {
			task.ClosedInterval = 1
		} else {
			task.ClosedInterval = 0
		}

		task.strStartId = curId
		task.strCurrId = curId
		if strings.Compare(curId, max) >= 0 {
			curId = max
		} else {
			if "" == dsp.FromWhere {
				//转义
				//curId = Stringescape(curId)
				xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= '%s' order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, priKey, curId, priKey, MyDataSizePerTask)
				//xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= %s order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, priKey, curId, priKey, MyDataSizePerTask)
			} else {
				//转义
				//curId = Stringescape(curId)
				xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE (%s) and `%s` >= '%s' order by `%s`  LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, dsp.FromWhere, priKey, curId, priKey, MyDataSizePerTask)
				//xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE (%s) and `%s` >= %s order by `%s`  LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, dsp.FromWhere, priKey, curId, priKey, MyDataSizePerTask)
			}

			LogDebug("Sql: %s", xSql)

			row := dbConn.QueryRow(xSql)
			err = row.Scan(&rowCur)
			if nil != err && sql.ErrNoRows != err {
				return errors.Trace(err)
			}
			if sql.ErrNoRows == err || strings.Compare(rowCur.String, max) >= 0 {
				curId = max
			} else {
				curId = rowCur.String
			}
		}

		task.strEndId = curId
		//转义
		//task.strEndId = Stringescape(curId)
		if dsp.DestTo == 0 && destEndpointSize > 1 {
			// round robin
			task.DestEndpoint = dsp.Dest.Endpoints[i%destEndpointSize]
		}

		LogDebug("Task: %v", task)
		taskChan.In() <- task
		atomic.AddInt32(&currState.SumTasks, 1)
	}

	return nil
}