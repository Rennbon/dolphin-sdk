package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gitlab.2se.com/hashhash/server-sdk/pb"
	"reflect"
)

var (
	delegate = &mpdelegate{
		direction: make(map[string]map[string]map[string]*grpcMethod),
	}

	ErrOverloadNotSupported = errors.New("The registered service does not support overloading of version,resource,action")
	ErrParamNotSpecified    = errors.New("Parameter not specified")
)

type grpcMethod struct {
	method reflect.Method
	argin  reflect.Type
	argout reflect.Type
}

type mpdelegate struct {
	//resource version action method
	direction map[string]map[string]map[string]*grpcMethod
}

func (m *mpdelegate) registerMethod(version, resource, action string, mehtod reflect.Method, int, out reflect.Type) error {
	if m.direction[resource] == nil {
		m.direction[resource] = make(map[string]map[string]*grpcMethod)
	}
	if m.direction[resource][version] == nil {
		m.direction[resource][version] = make(map[string]*grpcMethod)
	}
	if _, ok := m.direction[resource][version][action]; ok {
		return ErrOverloadNotSupported
	}
	m.direction[resource][version][action] = &grpcMethod{
		method: mehtod,
		argin:  int,
		argout: out,
	}
	return nil
}
func (m *mpdelegate) invoke(req *pb.ClientComRequest) *pb.ServerComResponse {
	response := &pb.ServerComResponse{
		Id: req.Id,
	}
	grpcM := m.direction[req.Meta.Resource][req.Meta.Revision][req.Meta.Action]
	inputs := make([]reflect.Value, 1)
	tmp := reflect.New(grpcM.argin).Interface().(descriptor.Message)
	err := ptypes.UnmarshalAny(req.Params, tmp)
	if err != nil {
		fmt.Println(err)
		response.Code = 400
		response.Text = ErrParamNotSpecified.Error()
		return response
	}
	inputs[0] = reflect.ValueOf(tmp)
	vals := grpcM.method.Func.Call(inputs)
	if vals[0].Type() == grpcM.argout && !vals[0].IsNil() {
		object, err := ptypes.MarshalAny(vals[0].Interface().(proto.Message))
		if err != nil {
			response.Code = 500
			response.Text = err.Error()
			return response
		}
		response.Code = 200
		response.Body = object
	}

	if !vals[1].IsNil() {
		response.Code = 500
		response.Text = vals[1].Interface().(error).Error()
	}
	return response
}