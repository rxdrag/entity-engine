package repository

import (
	"fmt"

	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
)

type QueryArg = map[string]interface{}

func Query(node graph.Noder, args QueryArg) []InsanceData {
	con, err := Open()
	if err != nil {
		panic(err.Error())
	}
	return con.doQueryNode(node, args)
}

func QueryOne(node graph.Noder, args QueryArg) interface{} {
	con, err := Open()
	if err != nil {
		panic(err.Error())
	}
	return con.doQueryOneNode(node, args)
}

func SaveOne(instance *data.Instance) (interface{}, error) {
	con, err := Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.BeginTx()
	defer con.ClearTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	obj, err := con.doSaveOne(instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return obj, nil
}

func InsertOne(instance *data.Instance) (interface{}, error) {
	con, err := Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer con.ClearTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	obj, err := con.doInsertOne(instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return obj, nil
}

func BatchQueryAssociations(association *graph.Association, ids []uint64) []map[string]interface{} {
	con, err := Open()
	if err != nil {
		panic(err.Error())
	}
	return con.doBatchAssociations(association, ids)
}

func IsEntityExists(name string) bool {
	con, err := Open()
	if err != nil {
		panic(err.Error())
	}
	return con.doCheckEntity(name)
}

func Install() error {
	con, err := Open()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = con.BeginTx()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	sql := `CREATE TABLE meta (
		id bigint NOT NULL AUTO_INCREMENT,
		content json DEFAULT NULL,
		publishedAt datetime DEFAULT NULL,
		createdAt datetime DEFAULT NULL,
		updatedAt datetime DEFAULT NULL,
		status varchar(45) DEFAULT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1507236403010867251 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	`
	_, err = con.Dbx.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	sql = `CREATE TABLE ability (
		id bigint NOT NULL AUTO_INCREMENT,
		entityUuid varchar(255) NOT NULL,
		columnUuid varchar(255) DEFAULT NULL,
		can tinyint(1) NOT NULL,
		expression text NOT NULL,
		abilityType varchar(128) NOT NULL,
		roleId bigint NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=4503621102206976 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	`
	_, err = con.Dbx.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	sql = `CREATE TABLE entity_auth_settings (
		id bigint NOT NULL AUTO_INCREMENT,
		entityUuid varchar(255) NOT NULL,
		expand tinyint(1) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY entityUuid_UNIQUE (entityUuid)
	) ENGINE=InnoDB AUTO_INCREMENT=4503616807239680 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	`
	_, err = con.Dbx.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = con.Commit()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
