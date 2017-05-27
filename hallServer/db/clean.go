package db

func Destory() {
	if BaseDB != nil && BaseDB.DB != nil {
		BaseDB.DB.Close()
	}

	if DB != nil && DB.DB != nil {
		DB.DB.Close()
	}

	if StatsDB != nil && StatsDB.DB != nil {
		StatsDB.DB.Close()
	}
}
