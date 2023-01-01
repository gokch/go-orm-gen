package db_mysql

import (
	"strings"

	"github.com/gokch/ornn/db"
)

func NewVendor(db *db.Conn) *Vendor {
	return &Vendor{
		db: db,
	}
}

type Vendor struct {
	db *db.Conn
}

func (t *Vendor) ConvType(dbType string) (genType string) {
	var unsigned bool

	opts := strings.Split(string(dbType), " ")
	for _, opt := range opts {
		opt = strings.ToLower(opt)
		if opt == "unsigned" {
			unsigned = true
		}
	}

	if len(opts) == 0 {
		return ""
	}

	fieldTypeWithLen := opts[0]
	pos := strings.Index(fieldTypeWithLen, "(")
	if pos != -1 {
		fieldTypeWithLen = fieldTypeWithLen[0:pos]
	}

	return t.convType(fieldTypeWithLen, unsigned)
}

func (t *Vendor) convType(dbType string, unsigned bool) string {
	switch strings.ToLower(dbType) {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
		return "string"
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return "[]byte"
	case "tinyint":
		if unsigned == true {
			return "uint8"
		}
		return "int8"
	case "smallint":
		if unsigned == true {
			return "uint16"
		}
		return "int16"
	case "int":
		if unsigned == true {
			return "uint32"
		}
		return "int32"
	case "bigint":
		if unsigned == true {
			return "uint64"
		}
		return "int64"
	case "float":
		return "float32"
	case "double", "real":
		return "float64"
	default:
		return "interface{}"
	}
}
