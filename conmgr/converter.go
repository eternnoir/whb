package conmgr

import (
	"errors"
	"fmt"

	"github.com/labstack/echo"
)

var DefuaultConverterMgr *ConverterMgr

type Converter interface {
	SourceName() string
	TargetName() string
	Process(e echo.Context) error
}

var (
	ErrConverterNotFound = errors.New("converter not found")
)

func init() {
	DefuaultConverterMgr = NewConverterMgr()
}

func NewConverterMgr() *ConverterMgr {
	return &ConverterMgr{
		convetMap: map[string]Converter{},
	}
}

type ConverterMgr struct {
	convetMap map[string]Converter
}

func (cm *ConverterMgr) RegConverter(con Converter) error {
	cm.convetMap[BuildSig(con.SourceName(), con.TargetName())] = con
	return nil
}

func (cm *ConverterMgr) Get(source, target string) (Converter, error) {
	sig := BuildSig(source, target)
	if val, ok := cm.convetMap[sig]; ok {
		return val, nil
	}
	return nil, ErrConverterNotFound
}
func BuildSig(source, target string) string {
	return fmt.Sprintf("%s:::%s", source, target)
}
