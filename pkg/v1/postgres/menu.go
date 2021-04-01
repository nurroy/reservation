package postgres

import (
	mnpb "belajar/reservation/proto/v1/menu"
	"database/sql"
	"fmt"
)

const (
	SELECT_MENU_LIMIT1 = `select * from reservation.menu limit 1`
	SELECT_MENU_LIST = `select * from reservation.menu`
	SELECT_MENU_BY_Name = `select * from reservation.menu where name = %s`
	INSERT_MENU = `insert into reservation.menu (name,"desc") values($1,$2)`
)

type Menu struct {
	ID string `db:"id"`
	Name string `db:"name"`
	Desc string	`db:"desc"`
}

func InsertMenu(data *mnpb.RetrieveRequest)error{
	_,err:= reservationDB.Exec(
		INSERT_MENU,data.GetMenu(),data.GetDescription())
	if err != nil {
		return err
	}
	return nil
}

func SelectMenu()([]*Menu, error){
	db, err := sql.Open(POSTGRES, pgConn)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(fmt.Sprintf(SELECT_MENU_LIST))

	defer db.Close()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rowsScanArr []*Menu

	//Fetch data to struct
	for rows.Next() {
		var rowsScan Menu
		err := rows.Scan(&rowsScan.ID,&rowsScan.Name,&rowsScan.Desc)

		if err != nil {
			return nil, err
		}

		// Append for every next row
		rowsScanArr = append(rowsScanArr, &rowsScan)
	}

	return rowsScanArr, nil
}
