package dal

import (
	"github.com/piotrjaromin/go-login-backend/web"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type entity struct {
	Id          string `bson:"_id"`
	Description string `bson:"description"`
}

type containerEntity struct {
	Id       string   `bson:"_id"`
	Entities []entity `bson:"entities"`
}

func TestRepository(t *testing.T) {

	repo := Create(DalConfig{
		Server:     "localhost",
		Database:   "testDB",
		Collection: "test",
	})

	repo.ClearAll()

	Convey("Repository should", t, func() {

		containerTestId := "containerTestId"
		testId := "testId"
		testEntity := entity{Id: testId, Description: "test"}

		Convey("save data", func() {

			isDup, err := repo.Save(testEntity)
			So(err, ShouldBeNil)
			So(isDup, ShouldBeFalse)
		})

		Convey("get data by id", func() {

			en := entity{}
			err := repo.GetById(testId, &en)

			So(err, ShouldBeNil)

			So(en.Id, ShouldEqual, testId)
			So(en.Description, ShouldEqual, testEntity.Description)
		})

		Convey("get all data", func() {

			ens := make([]entity, 0)
			err := repo.GetAll(&ens, web.DefaultPagination())

			So(err, ShouldBeNil)

			So(len(ens), ShouldEqual, 1)
			So(ens[0].Description, ShouldEqual, testEntity.Description)
			So(ens[0].Id, ShouldEqual, testId)
		})

		Convey("append data to array", func() {

			upsertErr := repo.Upsert(containerTestId, bson.M{
				"$push": bson.M{"entities": testEntity},
			})

			containerTestEntity := containerEntity{}
			getErr := repo.GetById(containerTestId, &containerTestEntity)

			So(getErr, ShouldBeNil)
			So(upsertErr, ShouldBeNil)

			So(len(containerTestEntity.Entities), ShouldEqual, 1)
			So(containerTestEntity.Entities[0].Description, ShouldEqual, testEntity.Description)
			So(containerTestEntity.Entities[0].Id, ShouldEqual, testId)
		})

		Convey("remove element from array", func() {

			updateErr := repo.Update(containerTestId, bson.M{
				"$pull": bson.M{"entities": bson.M{"_id": testEntity.Id}},
			})

			So(updateErr, ShouldBeNil)

			containerTestEntity := containerEntity{}
			getErr := repo.GetById(containerTestId, &containerTestEntity)

			So(getErr, ShouldBeNil)

			So(len(containerTestEntity.Entities), ShouldEqual, 0)
		})

		Convey("delete data by id", func() {

			query := NewQueryBuilder().WithField("_id", testId).Build()

			delErr := repo.DeleteByQuery(query)
			So(delErr, ShouldBeNil)

			en := entity{}
			err := repo.GetById(testId, &en)

			So(err, ShouldBeNil)
			So(en.Id, ShouldBeEmpty)
		})
	})
}
