package punct

import "fmt"

var RuntimeDLLs = []string{"VCRUNTIME140.dll", "VCRUNTIME140_1.dll", "MSVCP140.dll", "MSVCP140_1.dll"}

type MissingRuntime struct {
	File string
	Err  error
}

func (m *MissingRuntime) Error() string {
	return fmt.Sprintf("the Microsoft Visual C++ runtime is missing: %s cannot be loaded from the system folder (%v)", m.File, m.Err)
}

func CheckRuntime(load func(name string) error) error {
	for _, name := range RuntimeDLLs {
		if err := load(name); err != nil {
			return &MissingRuntime{File: name, Err: err}
		}
	}
	return nil
}
