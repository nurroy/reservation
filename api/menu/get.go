package menu

import (
	mnpb "belajar/reservation/proto/v1/menu"
	"context"
	pg "belajar/reservation/pkg/v1/postgres"
	"fmt"
)

func (s *Server)Get(ctx context.Context,request *mnpb.RetrieveRequest)(*mnpb.RetrieveResponse,error) {

	fmt.Println("MASUK")
	data,err := pg.SelectMenu()
	if err != nil{
		return nil,err
	}

	var rowsData []*mnpb.Data
	for i,_:=range data{
		row := &mnpb.Data{
			Menu:        data[i].Name,
			Description: data[i].Desc,
		}
		rowsData =append(rowsData,row)
	}

	return &mnpb.RetrieveResponse{Data: rowsData},nil

}
