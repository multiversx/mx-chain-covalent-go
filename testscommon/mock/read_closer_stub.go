package mock

type ReadCloserStub struct {
	ReadCalled  func(p []byte) (n int, err error)
	CloseCalled func() error
}

func (rcs *ReadCloserStub) Read(p []byte) (n int, err error) {
	if rcs.ReadCalled != nil {
		return rcs.ReadCalled(p)
	}

	return 0, nil
}

func (rcs *ReadCloserStub) Close() error {
	if rcs.CloseCalled != nil {
		return rcs.CloseCalled()
	}

	return nil
}
