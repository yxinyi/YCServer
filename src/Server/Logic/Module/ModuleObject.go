package module

import "time"

type ModuleInter interface {
	Init() error
	Start() error
	Stop() error
	Update(time_ time.Time)
}

type ModuleBase struct {
	m_name string
}

func (b *ModuleBase) Init() error            { return nil }
func (b *ModuleBase) Start() error           { return nil }
func (b *ModuleBase) Stop() error            { return nil }
func (b *ModuleBase) Name() string           { return b.m_name }
func (b *ModuleBase) Update(time_ time.Time) {}
