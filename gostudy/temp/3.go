//3、复合主键是整形时任务派发逻辑
func mysqlDispatchTaskByPriInt(dsp moverconfig.DataSourcePair, taskChan *channels.InfiniteChannel, fieldList []string, priKey string, priKeyTyp int) error {
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
	var rowSum sql.NullInt64
	var rowMin, rowMax, rowCount sql.NullInt64
	var xSql string
	if len(conditions) == 0 {
		xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s`", priKey, priKey, dsp.From.DBName, dsp.From.TableName)
	} else {
		xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s` WHERE %s", priKey, priKey, dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "))
	}

	var countSql string

	if len(conditions) == 0 {
		countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, priKey)
	} else {
		countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` WHERE %s group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "), priKey)
	}

	LogDebug("Sql: %s CountSql :%s", xSql, countSql)

	row := dbConn.QueryRow(xSql)
	err = row.Scan(&rowSum, &rowMin, &rowMax)
	if nil != err {
		logger.Errorf("mysqlDispatchTaskByInt-sql=%s-err=%v", xSql, err)
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
		logger.Errorf("mysqlDispatchTaskByInt-countSql=%s-err=%v", countSql, err)
		return errors.Trace(err)
	}

	var sum int64
	var min, max int64
	sum = rowSum.Int64
	if rowMin.Valid {
		min = rowMin.Int64
	}
	if rowMax.Valid {
		max = rowMax.Int64
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

	//限流，联合主键与rowCount.Int64有关
	if ((int(rowCount.Int64)+1)*len(fieldList))*2 < 65535 {
		task.RowCountLimit = int(rowCount.Int64) + 1
		if task.RowCountLimit < 100 {
			task.RowCountLimit = 100
		}
	} else {
		task.RowCountLimit = 65535 / (len(fieldList) * 2)
	}

	destEndpointSize := len(dsp.Dest.Endpoints)
	curId := min

	/*
		if curId >= max {
					LogWarn("Premature termination, curId=%d, maxId=%d, dsp: %v", curId, max, dsp)
					break
				}

				task.TaskID = int64(i)
				task.StartId = curId
				task.CurrId = curId
				curId += step
				if curId > max {
					curId = max
				}
				task.EndId = curId

				if dsp.DestTo == 0 && destEndpointSize > 1 {
					// round robin
					task.DestEndpoint = dsp.Dest.Endpoints[i%destEndpointSize]
				}

				LogDebug("Task: %+v", task)
				taskChan.In() <- task
				atomic.AddInt32(&currState.SumTasks, 1)
	*/

	for i := 0; i < taskSize; i++ {
		var rowCur sql.NullInt64
		if curId >= max {
			LogWarn("Premature termination, curId=%d, maxId=%d, dsp: %v", curId, max, dsp)
			break
		}
		task.TaskID = int64(i)
		task.StartId = curId
		task.CurrId = curId

		if curId >= max {
			curId = max
		} else {
			if "" == dsp.FromWhere {
				xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= %d order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, priKey, curId, priKey, MyDataSizePerTask)
			} else {
				xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE (%s) and `%s` >= %d  order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, dsp.FromWhere, priKey, curId, priKey, MyDataSizePerTask)
			}

			LogDebug("Sql: %s", xSql)

			row := dbConn.QueryRow(xSql)
			err = row.Scan(&rowCur)
			if nil != err && sql.ErrNoRows != err {
				return errors.Trace(err)
			}
			if sql.ErrNoRows == err || rowCur.Int64 >= max {
				curId = max
			} else {
				curId = rowCur.Int64
			}
		}

		task.EndId = curId
		appendTask(fmt.Sprintf("taskId=%d,startId=%v,taskEndId=%v,currId=%v,endId=%v,rowCur.Int64=%v,xSql=%v", task.TaskID, task.StartId, task.CurrId, task.EndId, curId, rowCur.Int64, xSql))

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