package ornn

import (
	"fmt"
	"os"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
)

type ORNN struct {
	db   *db.Conn
	conf *config.Config
}

func (t *ORNN) Init(db *db.Conn, conf *config.Config) {
	t.db = db
	t.conf = conf
}

func (t *ORNN) GenCode(path string) (err error) {
	if t.conf == nil {
		return fmt.Errorf("config is emtpy")
	}

	// gen code
	gen := &Gen{}
	code, err := gen.Gen(t.db, t.conf, path)
	if err != nil {
		return err
	}

	// write code to file
	err = os.WriteFile(path, []byte(code), 0700)
	if err != nil {
		return err
	}
	return nil
}
