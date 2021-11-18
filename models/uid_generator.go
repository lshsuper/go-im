package models

import (
	"github.com/gofrs/uuid"
	"strings"
)

type UIDGenerator struct {
}

func ( g UIDGenerator)NewID()string  {
     str,_:=uuid.NewV4()
	return  strings.Replace(str.String(),"-","",-1)
}
