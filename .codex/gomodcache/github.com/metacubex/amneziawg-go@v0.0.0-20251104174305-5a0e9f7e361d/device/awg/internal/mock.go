package internal

import (
	"github.com/metacubex/amneziawg-go/device/awg"
)

type mockGenerator struct {
	size int
}

func NewMockGenerator(size int) mockGenerator {
	return mockGenerator{size: size}
}

func (m mockGenerator) Generate(protocol *awg.Protocol) []byte {
	return make([]byte, m.size)
}

func (m mockGenerator) Size() int {
	return m.size
}

func (m mockGenerator) Name() string {
	return "mock"
}

type mockByteGenerator struct {
	data []byte
}

func NewMockByteGenerator(data []byte) mockByteGenerator {
	return mockByteGenerator{data: data}
}

func (bg mockByteGenerator) Generate(protocol *awg.Protocol) []byte {
	return bg.data
}

func (bg mockByteGenerator) Size() int {
	return len(bg.data)
}
