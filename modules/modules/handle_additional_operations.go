package modules

// RunAdditionalOperations implements AdditionalOperationsModule
func (m *Module) RunAdditionalOperations() error {
	return m.db.InsertEnableModules(m.cfg.Modules)
}
