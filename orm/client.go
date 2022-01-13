package orm

import (
	"google.golang.org/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/orm/model/ormschema"

	"github.com/cosmos/cosmos-sdk/orm/encoding/encodeutil"
	"github.com/cosmos/cosmos-sdk/orm/model/kvstore"
	"github.com/cosmos/cosmos-sdk/orm/model/ormlist"
	"github.com/cosmos/cosmos-sdk/orm/model/ormtable"
	"github.com/cosmos/cosmos-sdk/orm/types/ormerrors"
)

type ReadDBConnection struct {
	Schema      ormschema.Schema
	ReadBackend kvstore.ReadBackend
}

type DBConnection struct {
	*ReadDBConnection
	Backend kvstore.Backend
}

func (c ReadDBConnection) Get(message proto.Message, fieldNames ormtable.FieldNames, fields ...interface{}) (found bool, err error) {
	table, err := c.Schema.GetTable(message)
	if err != nil {
		return false, err
	}

	index := table.GetUniqueIndex(fieldNames)
	if index == nil {
		return false, ormerrors.CantFindIndex.Wrapf(
			"can't find unique index on table %s for fields %s",
			message.ProtoReflect().Descriptor().FullName(),
			fieldNames,
		)
	}

	return index.Get(c.ReadBackend, encodeutil.ValuesOf(fields...), message)
}

func (c ReadDBConnection) Has(message proto.Message, fieldNames ormtable.FieldNames, fields ...interface{}) (found bool, err error) {
	table, err := c.Schema.GetTable(message)
	if err != nil {
		return false, err
	}

	index := table.GetUniqueIndex(fieldNames)
	if index == nil {
		return false, ormerrors.CantFindIndex.Wrapf(
			"can't find unique index on table %s for fields %s",
			message.ProtoReflect().Descriptor().FullName(),
			fieldNames,
		)
	}

	return index.Has(c.ReadBackend, encodeutil.ValuesOf(fields...))
}

func (c DBConnection) Save(message proto.Message) error {
	table, err := c.Schema.GetTable(message)
	if err != nil {
		return err
	}

	return table.Save(c.Backend, message, ormtable.SAVE_MODE_DEFAULT)
}

func (c DBConnection) Insert(message proto.Message) error {
	table, err := c.Schema.GetTable(message)
	if err != nil {
		return err
	}

	return table.Save(c.Backend, message, ormtable.SAVE_MODE_INSERT)
}

func (c DBConnection) Update(message proto.Message) error {
	table, err := c.Schema.GetTable(message)
	if err != nil {
		return err
	}

	return table.Save(c.Backend, message, ormtable.SAVE_MODE_UPDATE)
}

func (c DBConnection) Delete(message proto.Message) error {
	table, err := c.Schema.GetTable(message)
	if err != nil {
		return err
	}

	return table.DeleteMessage(c.Backend, message)
}

func (c ReadDBConnection) List(message proto.Message, options ...ormlist.Option) (ormtable.Iterator, error) {
	table, err := c.Schema.GetTable(message)
	if err != nil {
		return nil, err
	}

	return ormlist.Iterator(c.ReadBackend, table, options...)
}