package main

database :=db.DB

tx,err :=database.Begin()
if err!=nil{
	return err
}

stmt ,err :=tx.Prepare(sqlQuery)
if err!=nil{
	tx.Rollback()
	return err
}
_,err :=stmt.Exec(paras...)
if err!=nil{
	tx.Rollback()
	return err
}

err=tx.Commit()
if err!=nil{
	tx.Rollback()
	return err
}

以上是我们使用事务时的一般操作，如果每做一次事务的操作均要进行重新写一遍代码岂不是很麻烦，尤其是出错时，Rollback需要多次在不同错误的地方的进行调用处理。

简单封装
第一步：采用defer处理Rollback

defer tx.Rollback()


无论成功与否，均进行Rollback操作，只是有点影响，如果成功还调用Rollback的话，将会报错。虽然可以忽略，但作为程序员，有必要进一步调整。



第二步：根据执行结果来选择执行Rollback，避免无效使用。

defer func() { //根据执行结果选择执行Rollback
    if err != nil && tx != nil {
    log.Println("ExecSqlWithTransaction defer err :", err)
        tx.Rollback()
    }
}()
如此，我们就可以根据事务的执行结果决定是否Rollback了。

第三步：封装
以上代码本身就具有极大的普适性，因此，我们抽出通用的参数，将此过程封装成一个func，以后就可以直接调用了。
func ExecSqlWithTransaction(database *sql.DB, query string, args ...interface{}) (err error) {
    tx, err := database.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil && tx != nil {
            tx.Rollback()
        }
    }()
    stmt, err := tx.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()
    _, err = stmt.Exec(args...)
    if err != nil {
        return err
    }
    return tx.Commit()
}

封装后我们可以如下使用：
if err := ExecSqlWithTransaction(database,sqlQuery,paras...);err != nil{
    //错误处理
}
封装后是不是很简洁啊？


进一步封装
在一个事务中可能会出现多个SELECT、UPDATE等操作，以上封装仅处理了一种操作，还不能满足我们的实际需求，因此需要更进一步封装。

func ExecSqlWithTransaction(db *sql.DB, handle func(tx *sql.Tx) error) (err error) {
 tx, err := db.Begin()
 if err != nil {
  return err
 }
 defer func() {
  if err != nil {
   tx.Rollback()
  }
 }()
 if err = handle(tx); err != nil {
  return err
 }
 return tx.Commit()
}

在handle func内可以直接使用事务tx进行增删改查。