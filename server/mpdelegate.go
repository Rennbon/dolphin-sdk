package server

import (
	"errors"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"reflect"
)

var (
	delegate = &mpdelegate{
		services:  make(map[string]reflect.Value),
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
	services  map[string]reflect.Value
	direction map[string]map[string]map[string]*grpcMethod
}

func (m *mpdelegate) registerService(resource string, value reflect.Value) {
	m.services[resource] = value
}
func (m *mpdelegate) registerMethod(version, resource, action string, mehtod reflect.Method, in, out reflect.Type) error {
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
		argin:  in,
		argout: out,
	}
	return nil
}
func (m *mpdelegate) invoke(req *pb.ClientComRequest) *pb.ServerComResponse {
	tc.push(req.TraceId) //trace save
	response := &pb.ServerComResponse{
		Id: req.Id,
	}
	grpcM := m.direction[req.MethodPath.Resource][req.MethodPath.Revision][req.MethodPath.Action]
	inputs := make([]reflect.Value, 2)
	tmp := reflect.New(grpcM.argin).Interface().(descriptor.Message)
	err := ptypes.UnmarshalAny(req.Params, tmp)
	if err != nil {
		response.Code = 400
		response.Text = ErrParamNotSpecified.Error()
		return response
	}
	inputs[0] = m.services[req.MethodPath.Resource]
	inputs[1] = reflect.ValueOf(tmp)
	vals := grpcM.method.Func.Call(inputs)

	tc.remove() //trace remove

	if vals[0].Elem().Type() == grpcM.argout && !vals[0].IsNil() {
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
