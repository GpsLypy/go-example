func mysqlReadDataByRange(task *Task) ([]RowData, error) {
	LogDebug("Do task: %v", task)

	dataSizeLimit := getConfig().DataSizeLimit

	if dataSizeLimit <= 100 {
		dataSizeLimit = 100
	} else if dataSizeLimit > task.RowCountLimit {
		//增加错误处理逻辑
		LogError(EC_Config_ERR, "Invalid config,DataSizelimit(%d) too big,It should be under(%d) ", dataSizeLimit, task.RowCountLimit)
		return nil, errors.New("Invalid config")
	}
	// get conn
	ep := task.FromEndpoint
	mysqlAddr := getConnstr(ep, "", false)
	dbConn, err := getMysqlDB(mysqlAddr)
	if nil != err {
		return nil, err
	}

	// create sql
	if FDT_INT == task.PrimaryKeyType {
		// create sql
		query := fmt.Sprintf("SELECT `%s`,%s FROM `%s`.`%s` WHERE `%s`>%d and `%s`<=%d",
			task.PrimaryKey, task.FromField, task.FromDBName, task.FromTable, task.PrimaryKey, task.CurrId, task.PrimaryKey, task.EndId)
		if task.IsSelfPri {
			query = fmt.Sprintf("SELECT `%s`,%s FROM `%s`.`%s` WHERE `%s`>= %d and `%s`<=%d",
				task.PrimaryKey, task.FromField, task.FromDBName, task.FromTable, task.PrimaryKey, task.CurrId, task.PrimaryKey, task.EndId)
		}
		if "" != task.FromWhere {
			query += fmt.Sprintf(" AND (%s)", task.FromWhere)
		}
		if task.IsSelfPri && task.CurrId == task.EndId {
			task.StartSkip = true
			log.Debugf("==task.taskId=%v======task.currId=%v===task.EndId=%v", task.TaskID, task.CurrId, task.EndId)
			query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d,%d", task.PrimaryKey, task.Skip, dataSizeLimit)
		} else {
			query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d,%d", task.PrimaryKey, task.WinSize, dataSizeLimit)
		}

		// query
		LogDebug("Query sql: %s", query)

		var rows *sql.Rows
		for i := 0; true; i++ {
			rows, err = dbConn.Query(query)
			if nil != err {
				sleep_time := rand.Intn(3000)
				LogError(EC_Ignore_ERR, "Retry, i=%d, sleep=%d,sql=%s, err: %v", i, sleep_time, query, err)
				time.Sleep(time.Duration(sleep_time) * time.Millisecond)
			} else {
				break
			}
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		colNames, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		curId := task.CurrId //记录当前任务开始的ID值，用来和任务读取数据结束后的ID值比较，从而确定偏移量WInsize
		var priKeyValue int64
		cols := make([]interface{}, len(colNames)-1)
		colPtrs := make([]interface{}, len(colNames))
		colPtrs[0] = &priKeyValue
		for i := 1; i < len(colNames); i++ {
			colPtrs[i] = &cols[i-1]
		}

		// get data
		results := make([]RowData, 0)
		for rows.Next() {
			err = rows.Scan(colPtrs...)
			if err != nil {
				return nil, err
			}

			var rd RowData
			rd.Row = make([]interface{}, len(cols))
			for i, v := range cols {
				if nil != v {
					// convert every thing to string
					rd.Row[i] = byteToString(v.([]uint8))
				} else {
					rd.Row[i] = v
				}
			}
			results = append(results, rd)
			task.CurrId = priKeyValue
		}

		//利用类滑动窗口原理记录偏移量
		if task.CurrId == curId {
			task.WinSize += len(results)
		} else {
			task.WinSize = 0
		}

		if task.IsSelfPri && task.StartSkip {
			task.Skip += len(results)
		}

		return results, nil
	} else {
		func mysqlReadDataByStrRange(task *Task) ([]RowData, error) {
			LogDebug("Do task: %v", task)
			dataSizeLimit := getConfig().DataSizeLimit
			if task.IsSelfPri {
				if dataSizeLimit <= 100 {
					dataSizeLimit = 100
				} else if dataSizeLimit > task.RowCountLimit {
					//增加错误处理逻辑
					LogError(EC_Config_ERR, "Invalid config,DataSizelimit(%d) too big,It should be under(%d) ", dataSizeLimit, task.RowCountLimit)
					return nil, errors.New("Invalid config")
				}
			}
		
			// get conn
			ep := task.FromEndpoint
			mysqlAddr := getConnstr(ep, "", false)
			dbConn, err := getMysqlDB(mysqlAddr)
			if nil != err {
				return nil, err
			}
		
			operator := ">"
			// the 1st's task is closed interval "[]", the others is semi-closed interval "(]"
			if 1 == task.ClosedInterval {
				operator = ">="
				task.ClosedInterval = 0
			}
			if task.IsSelfPri {
				operator = ">="
			}
			// create sql
			//SELECT `id` FROM `t1` WHERE `id` >= '"1"' and `id` <= '"9"'ORDER BY `id` LIMIT 5;
			query := fmt.Sprintf("SELECT `%s`,%s FROM `%s`.`%s` WHERE `%s` %s '%s' and `%s` <= '%s'",
				task.PrimaryKey, task.FromField, task.FromDBName, task.FromTable, task.PrimaryKey, operator, task.strCurrId, task.PrimaryKey, task.strEndId)
			if "" != task.FromWhere {
				query += fmt.Sprintf(" AND (%s)", task.FromWhere)
			}
			if task.IsSelfPri && task.strCurrId == task.strEndId {
				task.StartSkip = true
				log.Debugf("==task.taskId=%v======task.currId=%v===task.EndId=%v", task.TaskID, task.strCurrId, task.strEndId)
				query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d,%d", task.PrimaryKey, task.Skip, dataSizeLimit)
			} else {
				query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d, %d", task.PrimaryKey, task.WinSize, dataSizeLimit)
			}
		
			// query
			LogDebug("Query sql: %s", query)
		
			var rows *sql.Rows
			for i := 0; true; i++ {
				rows, err = dbConn.Query(query)
				if nil != err {
					sleep_time := rand.Intn(3000)
					LogError(EC_Ignore_ERR, "Retry, i=%d, sleep=%d,sql=%s, err: %v", i, sleep_time, query, err)
					time.Sleep(time.Duration(sleep_time) * time.Millisecond)
				} else {
					break
				}
			}
			if err != nil {
				return nil, err
			}
			defer rows.Close()
		
			colNames, err := rows.Columns()
			if err != nil {
				return nil, err
			}
			//局部偏移ID
			Localoffset := task.strCurrId
			batch := true
			var currentBatchId string
			var priKeyValue string
			cols := make([]interface{}, len(colNames)-1)
			colPtrs := make([]interface{}, len(colNames))
			colPtrs[0] = &priKeyValue
			for i := 1; i < len(colNames); i++ {
				colPtrs[i] = &cols[i-1]
			}
		
			// get data
			results := make([]RowData, 0)
			for rows.Next() {
				err = rows.Scan(colPtrs...)
				if err != nil {
					return nil, err
				}
				var rd RowData
				rd.Row = make([]interface{}, len(cols))
				for i, v := range cols {
					if nil != v {
						// convert every thing to string
						rd.Row[i] = byteToString(v.([]uint8))
					} else {
						rd.Row[i] = v
					}
				}
				results = append(results, rd)
				task.strCurrId = priKeyValue
				if batch == true {
					currentBatchId = priKeyValue
					batch = false
				}
				//task.strCurrId = Stringescape(priKeyValue)
			}
			if strings.Compare(task.strCurrId, Localoffset) == 0 && strings.Compare(task.strCurrId, currentBatchId) == 0 {
				task.WinSize += len(results)
			} else if strings.Compare(task.strCurrId, currentBatchId) == 0 {
				task.WinSize = len(results)
			} else {
				task.WinSize = 0
			}
			if task.IsSelfPri && task.StartSkip {
				task.Skip += len(results)
			}
		
			return results, nil
		}
	}

}

