package httphelper

import (
	"encoding/json"
	"net/http"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type ConditionFactory[T any] func(req *http.Request) (database.Condition[T], error)

type CrudHandler interface {
	GetHandler
	PostHandler
	PutHandler
	DeleteHandler
	WithModelName(s string) CrudHandler
}

type GetHandler interface {
	GetModelName() string
	Get(id string, resp http.ResponseWriter, req *http.Request)
	GetAll(resp http.ResponseWriter, req *http.Request)
}

type PostHandler interface {
	GetModelName() string
	Create(resp http.ResponseWriter, req *http.Request)
}

type PutHandler interface {
	GetModelName() string
	Update(id string, resp http.ResponseWriter, req *http.Request)
}

type DeleteHandler interface {
	GetModelName() string
	Delete(id string, resp http.ResponseWriter, req *http.Request)
}

type crudHelper[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string] struct {
	d        database.CrudHelper[T, MODEL, IDTYPE]
	idParser func(string) (IDTYPE, error)

	modelName        string
	conditionFactory ConditionFactory[T]
}

func NewCrudHelper[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string](d database.CrudHelper[T, MODEL, IDTYPE],
	idParser func(string) (IDTYPE, error), conditionFactory ConditionFactory[T]) CrudHandler {
	return crudHelper[T, MODEL, IDTYPE]{
		d:                d,
		idParser:         idParser,
		conditionFactory: conditionFactory,
	}
}

func (h crudHelper[T, MODEL, IDTYPE]) WithModelName(s string) CrudHandler {
	h.modelName = s
	return h
}

func (h crudHelper[T, MODEL, IDTYPE]) GetModelName() string {
	if h.modelName != "" {
		return h.modelName
	}
	return h.d.GetTableName()
}

func (h crudHelper[T, MODEL, IDTYPE]) Get(idStr string, resp http.ResponseWriter, req *http.Request) {
	id, err := h.idParser(idStr)
	if err != nil {
		SendError("In valid id for "+h.d.GetTableName(), err, resp)
		return
	}

	condition, err := h.conditionFactory(req)
	if err != nil {
		SendError("In valid condition for "+h.d.GetTableName(), err, resp)
		return
	}

	condition.And(condition.New().Set("id", database.ConditionOperationEqual, id))

	data, err := h.d.Get(req.Context(), req.URL.Query()["project"], condition)
	if err != nil {
		SendError("Error while GET "+h.d.GetTableName(), err, resp)
		return
	}

	if len(data) == 0 {
		Send404(resp)
		return
	}

	SendData(&data[0], resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) GetAll(resp http.ResponseWriter, req *http.Request) {
	condition, err := h.conditionFactory(req)
	if err != nil {
		SendError("In valid condition for "+h.d.GetTableName(), err, resp)
		return
	}

	data, err := h.d.Get(req.Context(), req.URL.Query()["project"], condition)
	if err != nil {
		SendError("Error while GET "+h.d.GetTableName(), err, resp)
		return
	}

	SendData(data, resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) Create(resp http.ResponseWriter, req *http.Request) {
	var body MODEL
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		SendError("In valid data for "+h.d.GetTableName(), err, resp)
		return
	}

	data, err := h.d.Create(req.Context(), &body)
	if err != nil {
		SendError("Error while POST "+h.d.GetTableName(), err, resp)
		return
	}

	SendData(data, resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) Update(idStr string, resp http.ResponseWriter, req *http.Request) {
	var body MODEL
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		SendError("In valid data for "+h.d.GetTableName(), err, resp)
		return
	}

	id, err := h.idParser(idStr)
	if err != nil {
		SendError("In valid id for "+h.d.GetTableName(), err, resp)
		return
	}

	condition, err := h.conditionFactory(req)
	if err != nil {
		SendError("In valid condition for "+h.d.GetTableName(), err, resp)
		return
	}

	condition.And(condition.New().Set("id", database.ConditionOperationEqual, id))

	err = h.d.Update(req.Context(), &body, req.URL.Query()["project"], condition.And(condition.New().Set("id", database.ConditionOperationEqual, id)))
	if err != nil {
		SendError("Error while PUT "+h.d.GetTableName(), err, resp)
		return
	}

	data, err := h.d.Get(req.Context(), req.URL.Query()["project"], condition.And(condition.New().Set("id", database.ConditionOperationEqual, id)))
	if err != nil {
		SendError("Error while GET "+h.d.GetTableName(), err, resp)
		return
	}

	if len(data) == 0 {
		SendError("No data found", nil, resp)
		return
	}

	SendData(&data[0], resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) Delete(idStr string, resp http.ResponseWriter, req *http.Request) {
	id, err := h.idParser(idStr)
	if err != nil {
		SendError("In valid id for "+h.d.GetTableName(), err, resp)
		return
	}

	condition, err := h.conditionFactory(req)
	if err != nil {
		SendError("In valid condition for "+h.d.GetTableName(), err, resp)
		return
	}

	condition.And(condition.New().Set("id", database.ConditionOperationEqual, id))

	err = h.d.Delete(req.Context(), condition)
	if err != nil {
		SendError("Error while DELETE "+h.d.GetTableName(), err, resp)
		return
	}

	SendData(nil, resp)
}
