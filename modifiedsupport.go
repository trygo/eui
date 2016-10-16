package eui

type ModifiedSupport struct {
	modified bool
}

func (this *ModifiedSupport) SetModified(b bool) {
	if this.modified != b {
		this.modified = b
	}
}

func (this *ModifiedSupport) IsModified() bool {
	return this.modified
}
