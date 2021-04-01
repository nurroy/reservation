package menu

import (
	pg "belajar/reservation/pkg/v1/postgres"
	mnpb "belajar/reservation/proto/v1/menu"
	"context"
)

func (s *Server) Post(ctx context.Context,req *mnpb.RetrieveRequest)(*mnpb.RetrieveResponse,error)  {

	err := pg.InsertMenu(req)
	if err !=nil{
		return nil, err
	}
	//init response
	message := &mnpb.Message{
		Message_TYPE:         "SU",
		MessageDesc:          "SUCCESS",
		MessageCode:          "0000",
	}
	var msg []*mnpb.Message
	msg = append(msg,message)

	status := &mnpb.Status{
		Message:              msg,
	}

	pagination := &mnpb.Pagination{
		HasMoreRecords:       "N",
		NumRecReturned:       "0",
	}

	header := &mnpb.Header{
		Status:               status,
		Pagination:           pagination,
	}
	return &mnpb.RetrieveResponse{Status: status,Header: header}, nil
}
