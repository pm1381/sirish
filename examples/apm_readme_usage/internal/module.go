package internal

import "context"

//go:generate sirish -t TestModule -tg=false

type TestModule interface {
	DoTest1(ctx context.Context, req DoTest1Request, span int) (string, error)
	DoTest2(ctx context.Context, req *DoTest2Request) (*DoTest2Response, error)
	DoNext() (ctx context.Context, res interface{}, err error)
}

type testModule struct {
}

func (t *testModule) DoNext() (ctx context.Context, res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *testModule) DoTest1(ctx context.Context, req DoTest1Request, span int) (string, error) {
	return req.id, nil
}

func (t *testModule) DoTest2(ctx context.Context, req *DoTest2Request) (*DoTest2Response, error) {
	return &DoTest2Response{res: ""}, nil
}

func NewTestHandler() TestModule {
	return &testModule{}
}

type DoTest1Request struct {
	id string
}

type DoTest2Request struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}

type DoTest2Response struct {
	res interface{}
}
