// Code generated by MockGen. DO NOT EDIT.
// Source: internal/nextactions/fetcher.go

// Package mock_nextactions is a generated GoMock package.
package mock_nextactions

import (
	gomock "github.com/golang/mock/gomock"
	trello "github.com/stevecshanks/next-actions-in-go/api/internal/trello"
	reflect "reflect"
)

// MockTrelloClient is a mock of TrelloClient interface
type MockTrelloClient struct {
	ctrl     *gomock.Controller
	recorder *MockTrelloClientMockRecorder
}

// MockTrelloClientMockRecorder is the mock recorder for MockTrelloClient
type MockTrelloClientMockRecorder struct {
	mock *MockTrelloClient
}

// NewMockTrelloClient creates a new mock instance
func NewMockTrelloClient(ctrl *gomock.Controller) *MockTrelloClient {
	mock := &MockTrelloClient{ctrl: ctrl}
	mock.recorder = &MockTrelloClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTrelloClient) EXPECT() *MockTrelloClientMockRecorder {
	return m.recorder
}

// OwnedCards mocks base method
func (m *MockTrelloClient) OwnedCards() ([]trello.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OwnedCards")
	ret0, _ := ret[0].([]trello.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OwnedCards indicates an expected call of OwnedCards
func (mr *MockTrelloClientMockRecorder) OwnedCards() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OwnedCards", reflect.TypeOf((*MockTrelloClient)(nil).OwnedCards))
}

// CardsOnList mocks base method
func (m *MockTrelloClient) CardsOnList(listID string) ([]trello.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CardsOnList", listID)
	ret0, _ := ret[0].([]trello.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CardsOnList indicates an expected call of CardsOnList
func (mr *MockTrelloClientMockRecorder) CardsOnList(listID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CardsOnList", reflect.TypeOf((*MockTrelloClient)(nil).CardsOnList), listID)
}

// ListsOnBoard mocks base method
func (m *MockTrelloClient) ListsOnBoard(boardID string) ([]trello.List, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListsOnBoard", boardID)
	ret0, _ := ret[0].([]trello.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListsOnBoard indicates an expected call of ListsOnBoard
func (mr *MockTrelloClientMockRecorder) ListsOnBoard(boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListsOnBoard", reflect.TypeOf((*MockTrelloClient)(nil).ListsOnBoard), boardID)
}
