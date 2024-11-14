func InitDB() {

	db, err := util.GetPgConn()
	if err != nil {
		logger.Info("Fail to connect pg:%v", err)
	}

	defer db.Close()

	exists := ensureSchemaExists(db, "ioit")

	if !exists {
		file, err := os.Open("./sql/ioit.sql")
		if err != nil {
			logger.Info("Fail to open file:%v", err)
		}
		defer file.Close()

		if err := executeSQLFile(db, file); err != nil {
			logger.Info("Fail to execute file:%v", err)
		}
	}

	schemaName := "ioit"
	tableNames := []interface{}{&vo.CreateJobRequest{}, &vo.JobNode{}}

	logger.Info("start init job_info")
	newDb := util.GlobalGormDb.Exec(fmt.Sprintf("SET search_path TO %s", schemaName))

	for _, table := range tableNames {
		if !newDb.Migrator().HasTable(table) {
			err := newDb.Migrator().CreateTable(table)
			if err != nil {
				logger.Info("failed to create table:%v", err)
			}
		}
	}

	var constraints string
	err = util.GlobalGormDb.Raw(`SELECT conname FROM pg_constraint WHERE conrelid = 'ioit.job_info'::regclass AND contype = 'p'`).Scan(&constraints).Error
	if err != nil {
		logger.Info("select pkey error:%v", err)
	}
	logger.Info("constraints:%v", constraints)
	if constraints != "" {
		err = util.GlobalGormDb.Exec(fmt.Sprintf(`ALTER TABLE ioit.job_info DROP CONSTRAINT %s`, constraints)).Error
		if err != nil {
			logger.Info("select pkey error:%v", err)

		}
	}

	var count int64
	err = util.GlobalGormDb.Raw(`SELECT COUNT(1) FROM pg_indexes WHERE schemaname = 'ioit' AND tablename = 'job_info' AND indexname = 'idx_namespace'`).Scan(&count).Error

	if err != nil {
		logger.Info("select index failed:%v", err)
	}

	if count == 0 {
		err := util.GlobalGormDb.Exec(`CREATE INDEX idx_namespace on ioit.job_info (namespace);`)
		if err != nil {
			logger.Info("create index failed:%v", err)
		}
	} else {
		logger.Info("index already exist")
	}

	addRetry()

}

func addRetry() {
	var count int64
	if err := util.GlobalGormDb.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = 'database_info' AND column_name = 'retry'").Scan(&count).Error; err != nil {
		logger.Info("Fail to check retry:%v", err)
	}

	if count == 0 {
		if err := util.GlobalGormDb.Migrator().AddColumn(&vo.DatabaseInfoRequest{}, "Retry"); err != nil {
			logger.Info("Fail to add retry:%v", err)
		}

		if err := util.GlobalGormDb.Model(&vo.DatabaseInfoRequest{}).Where("retry IS NULL").Update("retry", "5").Error; err != nil {
			logger.Info("Fail to update retry:%v", err)
		}
	} else {
		logger.Info("retry already exist")
	}
}

func ensureSchemaExists(db *sql.DB, schemaName string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_namespace WHERE nspname = $1)`
	if err := db.QueryRow(query, schemaName).Scan(&exists); err != nil {
		logger.Info("Fail to check schema:%v", err)
	}

	if !exists {
		createSchemaSQL := fmt.Sprintf(`CREATE SCHEMA %s`, schemaName)
		if _, err := db.Exec(createSchemaSQL); err != nil {
			logger.Info("Fail to create schema:%v", err)
		}
	} else {
		logger.Info("scheam is exists")
	}
	return exists
}

func executeSQLFile(db *sql.DB, file io.Reader) error {
	reader := bufio.NewReader(file)
	var sb strings.Builder

	for {
		line, err := reader.ReadString(';')
		if err != nil && err != io.EOF {
			logger.Info("Fail to read sql:%v", err)
			return err
		}

		line = strings.TrimSpace(line)
		if line != "" {
			sb.WriteString(line)
		}
		if err == io.EOF {
			break
		}
	}

	queries := strings.Split(sb.String(), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query != "" {
			if _, err := db.Exec(query); err != nil {
				logger.Info("Fail to excute sql:%v", err)
				return err
			}
		}
	}
	return nil
}
