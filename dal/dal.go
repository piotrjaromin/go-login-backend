package dal

import (
	"github.com/op/go-logging"
	"github.com/piotrjaromin/go-login-backend/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DalConfig struct {
	Server     string
	Database   string
	Collection string
}

type Dal struct {
	Save            func(interface{}) (bool, error)
	GetById         func(id string, entity interface{}) error
	GetAll          func(container interface{}, pagination web.Pagination) error
	GetByQuery      func(container interface{}, pagination web.Pagination, query Query) error
	Upsert          func(id string, element interface{}) error
	Update          func(id string, element interface{}) error
	UpdateByQuery   func(query Query, element interface{}) error
	ClearAll        func() error
	DeleteByQuery   func(query Query) error
	DeleteById      func(id string) error
	AddToArray      func(id string, field string, element interface{}) error
	DeleteFromArray func(id string, query Query) error
}

func Create(repoConfig DalConfig) Dal {
	var log = logging.MustGetLogger("[GenericDal]")

	log.Info("Creating repository for ", repoConfig.Collection)
	//mgo.SetLogger(log.New(os.Stdout, "[MGO] ", 1))

	session, err := mgo.Dial(repoConfig.Server)

	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(repoConfig.Database).C(repoConfig.Collection)

	getById := func(id string, entity interface{}) error {

		err = c.Find(bson.M{"_id": id}).One(entity)
		if err == nil || err.Error() == "not found" {
			return nil
		}

		log.Errorf("[MGO] getById error, details: ", err)
		return err
	}

	deleteByQuery := func(query Query) error {

		return c.Remove(query.fields)
	}

	getByQuery := func(container interface{}, pagination web.Pagination, query Query) error {

		err := c.Find(query.fields). //
						Sort(query.sort...). 					 //
						Skip((pagination.PageNumber - 1) * pagination.PageSize). //
						Limit(pagination.PageSize).                              //
						All(container)

		return err
	}

	getAll := func(container interface{}, pagination web.Pagination) error {

		return getByQuery(container, pagination, NewQueryBuilder().Build())
	}

	save := func(element interface{}) (bool, error) {

		if err := c.Insert(element); err != nil {
			if mgo.IsDup(err) {
				return true, nil
			}
			return false, err
		}

		return false, nil
	}

	upsert := func(id string, element interface{}) error {

		_, err := c.UpsertId(id, element)
		return err
	}

	update := func(id string, element interface{}) error {

		return c.UpdateId(id, element)
	}

	updateByQuery := func(query Query, element interface{}) error {

		return c.Update(query.fields, element)
	}

	clearAll := func() error {

		return c.DropCollection()
	}

	deleteById := func(id string) error {
		return c.Remove(bson.M{"_id": id})
	}

	addToArray := func(id string, field string, element interface{}) error {
		_, err := c.UpsertId(id, bson.M{
			"$push": bson.M{field: element},
		})

		return err
	}

	deleteFromArray := func(id string, query Query) error {
		_, err := c.UpsertId(id, bson.M{
			"$pull": query.fields,
		})

		return err
	}

	return Dal{
		Save:            save,
		GetById:         getById,
		GetAll:          getAll,
		GetByQuery:      getByQuery,
		Upsert:          upsert,
		ClearAll:        clearAll,
		DeleteByQuery:   deleteByQuery,
		Update:          update,
		UpdateByQuery:   updateByQuery,
		DeleteById:      deleteById,
		AddToArray:      addToArray,
		DeleteFromArray: deleteFromArray,
	}
}
