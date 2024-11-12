package jobschedule

// type Executor struct {
// 	Address string
// 	Weight  int
// }

// type Job struct {
// 	Id          string `json:"id"`
// 	Name        string `json:"name"`
// 	Schedule    string `json:"schedule"`
// 	Code        string `json:"code"`
// 	CurrentNode string `json:"s_node"`
// }

// type Schedule struct {
// 	executors []Executor
// 	jobs      []Job
// 	mu        sync.Mutex
// }

// func NewScheduler() *Schedule {
// 	return &Schedule{}
// }

// func (s *Schedule) AddExecutor(executor Executor) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.executors = append(s.executors, executor)
// }

// func (s *Schedule) AddJob(job Job) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.jobs = append(s.jobs, job)
// }

// func (s *Schedule) selectExecutor() *Executor {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	db, err := GetPgConn()
// 	if err != nil {
// 		return nil
// 	}
// 	defer db.Close()

// 	rows, err := db.Query("SELECT addr, weight FROM hsm_scheduling.hosts")
// 	if err != nil {
// 		logger.Info("Failed to select addr")
// 		return nil
// 	}
// 	defer rows.Close()

// 	s.executors = []Executor{}

// 	for rows.Next() {
// 		var executor Executor
// 		if err := rows.Scan(&executor.Address, &executor.Weight); err != nil {
// 			logger.Info("Failed to scan addr")
// 			continue
// 		}
// 		s.executors = append(s.executors, executor)
// 	}

// 	if err = rows.Err(); err != nil {
// 		logger.Info("Rows iteration failed:%v", err)
// 		return nil
// 	}

// 	if len(s.executors) == 0 {
// 		return nil
// 	}
// 	minWeight := s.executors[0].Weight
// 	selected := &s.executors[0]
// 	for i := 1; i < len(s.executors); i++ {
// 		if s.executors[i].Weight < minWeight {
// 			minWeight = s.executors[i].Weight
// 			selected = &s.executors[i]
// 		}
// 	}
// 	return selected
// }

// func (s *Schedule) Start(ctx context.Context) {
// 	// 传递过来一个任务结构体
// 	// 将任务添至调度中
// 	for _, job := range s.jobs {
// 		go func(job Job) {
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case <-time.After(time.Until(NextExecutionTime(job.Schedule))):
// 					executor := s.selectExecutor()
// 					if executor != nil {
// 						logger.Info("Assigning job to executor at %s", executor.Address)

// 						if err := s.executeJob(executor, job); err != nil {
// 							logger.Info("Failed to execute")
// 						}
// 						// 更新host表格的runningtask字段
// 						// UpdateRunningTask()
// 						// 更新权重
// 						// UpdateWeight()
// 					}
// 				}
// 			}
// 		}(job)
// 	}
// }

// func NextExecutionTime(schedule string) time.Time {
// 	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
// 	scheduleParsed, err := parser.Parse(schedule)
// 	if err != nil {
// 		logger.Info("cronexpr parse fail:%v", err)
// 		return time.Time{}
// 	}
// 	fmt.Printf("**********scheduleParsed: %v\n", scheduleParsed)
// 	return scheduleParsed.Next(time.Now())
// }

// func (s *Schedule) executeJob(executor *Executor, job Job) error {
// 	requestBody, err := json.Marshal(job)
// 	if err != nil {
// 		logger.Info("Fail to marshal:%v", err)
// 	}
// 	req, err := http.NewRequest("POST", executor.Address+"/api/hsm-io-it/execute", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		logger.Info("Fail to send request")
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	client := &http.Client{Timeout: 60 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		logger.Info("Failed to excute job")
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		logger.Info("Executor returned non-OK")
// 	}
// 	logger.Info("Successfully executed")
// 	return nil
// }
