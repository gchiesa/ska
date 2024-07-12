package part

import "os"

func (p *Part) WithContent(content []byte) *Part {
	p.content = content
	return p
}

func (p *Part) SetFilePath(path string) *Part {
	p.refFilePath = path
	return p
}

func (p *Part) CreateFile() error {
	return os.WriteFile(p.refFilePath, p.content, 0644)
}
