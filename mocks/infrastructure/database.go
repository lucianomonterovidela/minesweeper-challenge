package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBaseProviderMock struct {
	mock.Mock
}

func (m *DataBaseProviderMock) ConnectDatabase() (interface{}, error) {
	args := m.Called()
	err := args.Error(1)
	if args.Get(0) == nil {
		return "", err
	}
	return args.Get(0), args.Error(1)
}

func (m *DataBaseProviderMock) Insert(collectionName string, val interface{}) (interface{}, error) {
	args := m.Called(collectionName, val)
	err := args.Error(1)
	if args.Get(0) == nil {
		return nil, err
	}
	return args.Get(0), args.Error(1)
}

func (m *DataBaseProviderMock) Upsert(collectionName string, id interface{}, val interface{}) error {
	args := m.Called(collectionName, id, val)
	return args.Error(0)
}

func (m *DataBaseProviderMock) ReplaceById(collectionName string, id interface{}, val interface{}) error {
	args := m.Called(collectionName, id, val)
	return args.Error(0)
}

func (m *DataBaseProviderMock) GetById(collectionName string, id interface{}) (*mongo.SingleResult, error) {
	args := m.Called(collectionName, id)
	err := args.Error(1)
	if args.Get(0) == nil {
		return nil, err
	}
	return nil, args.Error(1)
}

func (m *DataBaseProviderMock) Find(collectionName string, filter interface{}, ops *options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(collectionName, filter)
	err := args.Error(1)
	if args.Get(0) == nil {
		return nil, err
	}
	return nil, args.Error(1)
}

func (m *DataBaseProviderMock) Aggregate(collectionName string, pipeline interface{}, opts *options.AggregateOptions) (*mongo.Cursor, error) {
	args := m.Called(collectionName, pipeline, opts)
	err := args.Error(1)
	if args.Get(0) == nil {
		return nil, err
	}
	return nil, args.Error(1)
}
