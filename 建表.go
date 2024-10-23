func MigrateTables(db *gorm.DB) error {
	// Migrate the U8TaskCenter table
	if err := db.AutoMigrate(&U8TaskCenter{}); err != nil {
		return err
	}

	// Migrate the U8Department table
	if err := db.AutoMigrate(&U8Department{}); err != nil {
		return err
	}

	// Migrate the U8Supplier table
	if err := db.AutoMigrate(&U8Supplier{}); err != nil {
		return err
	}

	// Migrate the U8Unit table
	if err := db.AutoMigrate(&U8Unit{}); err != nil {
		return err
	}

	// Migrate the U8CustomerData table
	if err := db.AutoMigrate(&U8CustomerData{}); err != nil {
		return err
	}

	// Migrate the U8Person table
	if err := db.AutoMigrate(&U8Person{}); err != nil {
		return err
	}

	// Migrate the U8ProductionOrder table
	if err := db.AutoMigrate(&U8ProductionOrder{}); err != nil {
		return err
	}

	// Migrate the U8Material table
	if err := db.AutoMigrate(&U8Material{}); err != nil {
		return err
	}

	// Migrate the U8MaterialBOM table
	if err := db.AutoMigrate(&U8MaterialBOM{}); err != nil {
		return err
	}

	// Migrate the KingdeeDepartment table
	if err := db.AutoMigrate(&KingdeeDepartment{}); err != nil {
		return err
	}

	// Migrate the KingdeeWorkCenter table
	if err := db.AutoMigrate(&KingdeeWorkCenter{}); err != nil {
		return err
	}

	// Migrate the KingdeeSupplier table
	if err := db.AutoMigrate(&KingdeeSupplier{}); err != nil {
		return err
	}

	// Migrate the KingdeePlanOrder table
	if err := db.AutoMigrate(&KingdeePlanOrder{}); err != nil {
		return err
	}

	// Migrate the KingdeeUnit table
	if err := db.AutoMigrate(&KingdeeUnit{}); err != nil {
		return err
	}

	// Migrate the KingdeeCustomer table
	if err := db.AutoMigrate(&KingdeeCustomer{}); err != nil {
		return err
	}

	// Migrate the KingdeeProductionOrder table
	if err := db.AutoMigrate(&KingdeeProductionOrder{}); err != nil {
		return err
	}

	// Migrate the KingdeeMaterial table
	if err := db.AutoMigrate(&KingdeeMaterial{}); err != nil {
		return err
	}

	// Migrate the KingdeeBOM table
	if err := db.AutoMigrate(&KingdeeBOM{}); err != nil {
		return err
	}

	// Migrate the KingdeeEmployee table
	if err := db.AutoMigrate(&KingdeeEmployee{}); err != nil {
		return err
	}

	return nil
}