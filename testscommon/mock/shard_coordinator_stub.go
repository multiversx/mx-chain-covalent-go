package mock

// ShardCoordinatorMock -
type ShardCoordinatorMock struct {
	SelfID          uint32
	ComputeIdCalled func(address []byte) uint32
}

// ComputeId -
func (scm *ShardCoordinatorMock) ComputeId(address []byte) uint32 {
	if scm.ComputeIdCalled != nil {
		return scm.ComputeIdCalled(address)
	}

	return 0
}

// SelfId -
func (scm *ShardCoordinatorMock) SelfId() uint32 {
	return scm.SelfID
}

// IsInterfaceNil returns true if there is no value under the interface
func (scm *ShardCoordinatorMock) IsInterfaceNil() bool {
	return scm == nil
}
