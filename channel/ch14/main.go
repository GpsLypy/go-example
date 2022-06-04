
func worker() error {
	for {
		select {
		case data := <-taskChan.Out():
			//此分支被选中，说明拿到了任务
			// do task
			// set state
		case <-ticker.C:
			//此分支被选中，检查任务中心中的任务是否全部做完，做完就退出
			/* finish dump worker */
			if moverconfig.GetCurrState().Prepared == true && taskChan.len() == 0 {
				return nil
			}
		case <-StopJobChan:
			logging.LogInfo("Worker(%d) was stop", wid)
			return nil

		case <-stopWorkerCh:
			logging.LogInfo("Worker(%d) stop", wid)
			return nil
		}
	}
}
