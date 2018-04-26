package converters

import "github.com/eternnoir/whb/conmgr"

func init() {
	conmgr.DefuaultConverterMgr.RegConverter(NewJenkins2HC())
	conmgr.DefuaultConverterMgr.RegConverter(NewCrashlytics2HC())
}
