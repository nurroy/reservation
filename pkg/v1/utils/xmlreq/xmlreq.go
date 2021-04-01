package xmlreq

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type (
	// XMLReq struct model to generate XML request header
	XMLReq struct {
		ServiceID        string
		ServiceReqVer    string
		ChannelID        string
		LangID           string
		OriginInstanceID string
		OriginVersion    string
		SchemaLocation   string
		SchemaInstance   string
		useGlobal        bool
	}

	// RequestHeader ...
	RequestHeader struct {
		MessageKey         MessageKey  `xml:"MessageKey"`
		Security           interface{} `xml:"Security"`
		RequestMessageInfo interface{} `xml:"RequestMessageInfo"`
	}

	// MessageKey ...
	MessageKey struct {
		GlobalUUID            string `xml:"GlobalUUID,omitempty"`
		RequestUUID           string `xml:"RequestUUID"`
		ServiceID             string `xml:"ServiceId"`
		ServiceRequestVersion string `xml:"ServiceRequestVersion"`
		ChannelID             string `xml:"ChannelId"`
		OriginatorInstanceID  string `xml:"OriginatorInstanceId,omitempty"`
		OriginatorVersion     string `xml:"OriginatorVersion,omitempty"`
		LanguageID            string `xml:"LanguageId"`
	}

	// FIXML ...
	FIXML struct {
		SchemaLocation string        `xml:"xsi:schemaLocation,attr"`
		SchemaInstance string        `xml:"xmlns:xsi,attr"`
		XMLNS          string        `xml:"xmlns,attr"`
		RequestHeader  RequestHeader `xml:"Header>RequestHeader"`
		RequestBody    interface{}   `xml:"Body"`
	}
)

var (
	errMissingRequired       = errors.New("missing required data, please check ServiceID/ChannelID/ServiceReqVer input")
	errMissingSchemaInstance = errors.New("missing schema instance")
	errMissingSchemaLoc      = errors.New("missing schema location")
	errMissingSecurity       = errors.New("missing security struct")
	errMissingMessage        = errors.New("missing message struct")
)

// NewXMLReq initiate object for generating xml request
func NewXMLReq(req XMLReq, useGlobalUUID bool) (*XMLReq, error) {
	if req.SchemaInstance == "" {
		return nil, errMissingSchemaInstance
	}
	if req.SchemaLocation == "" {
		return nil, errMissingSchemaLoc
	}
	if req.ServiceID == "" || req.ServiceReqVer == "" || req.ChannelID == "" {
		return nil, errMissingRequired
	}
	req.useGlobal = useGlobalUUID
	return &req, nil
}

// GenerateXMLRequest xml generate request, the request data should be in struct with correct XML schema, otherwise the request will be failed
func (x *XMLReq) GenerateXMLRequest(reqSecurity, reqMessage, reqBody interface{}) ([]byte, error) {
	var globalUID string
	uid := uuid.New().String()

	if reqSecurity != nil && reflect.TypeOf(reqSecurity).Kind() != reflect.Struct {
		return nil, errMissingSecurity
	}

	if reqMessage != nil && reflect.TypeOf(reqMessage).Kind() != reflect.Struct {
		return nil, errMissingMessage
	}

	// TODO: validate the interface struct, if its not contain the XML tag return as error

	if x.useGlobal {
		globalUID = uuid.New().String()
	}
	r := FIXML{
		XMLNS:          x.SchemaLocation,
		SchemaInstance: x.SchemaInstance,
		SchemaLocation: fmt.Sprintf("%s %s.xsd", x.SchemaLocation, x.ServiceID),
		RequestHeader: RequestHeader{
			MessageKey: MessageKey{
				GlobalUUID:            globalUID,
				RequestUUID:           uid,
				ServiceID:             x.ServiceID,
				ChannelID:             x.ChannelID,
				ServiceRequestVersion: x.ServiceReqVer,
				OriginatorInstanceID:  x.OriginInstanceID,
				OriginatorVersion:     x.OriginVersion,
				LanguageID:            x.LangID,
			},
			Security:           reqSecurity,
			RequestMessageInfo: reqMessage,
		},
		RequestBody: reqBody,
	}
	out, err := xml.MarshalIndent(r, "  ", "    ")
	if err != nil {
		return nil, err
	}

	return out, nil
}

// ParseXMLResponse is the method to parse the XML response to general struct response (finacle response)
func ParseXMLResponse(data []byte, target interface{}) (err error) {
	err = xml.Unmarshal(data, target)
	if err != nil {
		return err
	}
	return
}
