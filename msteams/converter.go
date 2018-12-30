package msteams

import "github.com/eternnoir/whb/conmgr"

func init() {
	conmgr.DefuaultConverterMgr.RegConverter(NewLine2Teams())
}
