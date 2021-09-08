package module

import "time"

var mgr = newModuleManager()

type ModuleManager struct {
	m_module_list map[string]ModuleInter
}

func newModuleManager() *ModuleManager {
	return &ModuleManager{
		m_module_list: make(map[string]ModuleInter),
	}
}

func (b *ModuleManager) register(module_name_ string, module_obj_ ModuleInter) error {
	_, exists := b.m_module_list[module_name_]
	if exists {
		panic("exists before " + module_name_)

	}
	b.m_module_list[module_name_] = module_obj_
	return nil
}

func (b *ModuleManager) init() error {
	for _, it := range b.m_module_list {
		_err := it.Init()
		if _err != nil {
			return _err
		}
	}
	return nil
}
func (b *ModuleManager) start() error {
	for _, it := range b.m_module_list {
		_err := it.Start()
		if _err != nil {
			return _err
		}
	}
	return nil
}
func (b *ModuleManager) stop() error {
	for _, it := range b.m_module_list {
		_err := it.Stop()
		if _err != nil {
			return _err
		}
	}
	return nil
}

func (b *ModuleManager) update(time_ time.Time) {
	for _, it := range b.m_module_list {
		it.Update(time_)
	}
}

func Register(module_name_ string, module_obj_ ModuleInter) error {
	return mgr.register(module_name_, module_obj_)
}
func Init() error {
	return mgr.init()
}
func Start() error {
	return mgr.start()
}
func Stop() error {
	return mgr.stop()
}
func Update(time_ time.Time) {
	mgr.update(time_)
}
