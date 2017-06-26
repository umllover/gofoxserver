package db

var userUpdateSql = [][]string{
<<<<<<< HEAD
	0 : []string{
	},
	1 : []string{
		"ALTER TABLE version_update_test ADD test11 int",
		"ALTER TABLE version_update_test ADD test12 int",
		"ALTER TABLE version_update_test ADD test13 int",
	},
	2 : []string{
		"ALTER TABLE version_update_test ADD test21 int",
		"ALTER TABLE version_update_test ADD test22 int",
		"ALTER TABLE version_update_test ADD test23 int",
	},
	3 : []string{
		"ALTER TABLE version_update_test ADD test31 int",
		"ALTER TABLE version_update_test ADD test32 int",
		"ALTER TABLE version_update_test ADD test33 int",
	},
	 //[0] ,"ALTER TABLE table_name ADD field_name field_type",
=======
	0: []string{
		"ALTER TABLE userattr ADD XXXX int;"
	},
>>>>>>> tbnn

	1: []string{
		"ALTER TABLE userattr ADD XXXX int;"
	},

	3: []string{
		"ALTER TABLE userattr ADD XXXX int;"
	},

	4: []string{
		"ALTER TABLE userattr ADD XXXX int;"
	},
}

///////////////////////////////////////////////////// log db /////////////////////////////////////////////////
var statsUpdateSql = [][]string{
	0 : []string{
	},
	1 : []string{
		"ALTER TABLE version_update_test ADD test11 int",
		"ALTER TABLE version_update_test ADD test12 int",
		"ALTER TABLE version_update_test ADD test13 int",
	},
	2 : []string{
		"ALTER TABLE version_update_test ADD test21 int",
		"ALTER TABLE version_update_test ADD test22 int",
		"ALTER TABLE version_update_test ADD test23 int",
	},
	3 : []string{
		"ALTER TABLE version_update_test ADD test31 int",
		"ALTER TABLE version_update_test ADD test32 int",
		"ALTER TABLE version_update_test ADD test33 int",
	},
}
