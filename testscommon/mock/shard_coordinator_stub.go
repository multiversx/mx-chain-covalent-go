package mock

// ShardCoordinatorMock -
type ShardCoordinatorMock struct {
	SelfID          uint32
	ComputeIdCalled func(address []byte) uint32
}

// ComputeId calls a custom compute id function if defined, otherwise returns 0
func (scm *ShardCoordinatorMock) ComputeId(address []byte) uint32 {
	if scm.ComputeIdCalled != nil {
		return scm.ComputeIdCalled(address)
	}

	return 0
}

// SelfId returns SelfID member
func (scm *ShardCoordinatorMock) SelfId() uint32 {
	return scm.SelfID
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (scm *ShardCoordinatorMock) IsInterfaceNil() bool {
	return scm == nil
}
