package gen

import (
	"database/sql"

	"github.com/gokch/ornn/db"
)

func NewLogic(db *db.Conn) *Logic {
	return &Logic{
		db: db,
	}
}

type Logic struct {
	db *db.Conn
}

func (t *Logic) NoTxExample() error {
	job := t.db.Job()

	schema := &Schema{}
	schema.Init(job)

	//---------------------- start ----------------------//
	// Execute generated function down here!!
	// _, err = gen.TestGroup.Select()

	//----------------------- end -----------------------//
	return nil
}

func (t *Logic) TxExample() error {
	job, err := t.db.TxJob(sql.LevelSerializable, false)
	if err != nil {
		return err
	}

	schema := &Schema{}
	schema.Init(job)

	//---------------------- start ----------------------//
	// Execute generated function down here!!
	// _, err = schema.TestGroup.Select()

	//----------------------- end -----------------------//
	if err != nil {
		job.Rollback()
		return err
	}
	err = job.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (t *Logic) TxFuncExample() error {
	return t.db.TxJobFunc(sql.LevelSerializable, false, func(job *db.Job) error {
		schema := &Schema{}
		schema.Init(job)

		//---------------------- start ----------------------//
		// Execute generated function down here!!
		// _, err = schema.TestGroup.Select()

		//----------------------- end -----------------------//

		return nil
	})
}
