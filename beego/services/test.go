package services

import (
	"StandardProject/beego/models"
	"StandardProject/common/gormdb"
	"context"
	"database/sql"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
)

type testService struct {
	appCtx context.Context
}

func NewTestService(appCtx context.Context) *testService {
	return &testService{appCtx: appCtx}
}

func (c *testService) Test1() (string, error) {
	studentModel := models.NewStudentModel(c.appCtx)
	studentModel.Test(models.Student{})
	return "vsasv", nil
}

func (c *testService) Test2(id int, name string) ([]models.Student, error) {
	studentModel := models.NewStudentModel(c.appCtx)
	//return studentModel.Find()
	where := models.Student{
		Id:   id,
		Name: name,
	}
	//return studentModel.SelectFieldsFindByStudent(nil, where)
	//studentModel.DeleteByPrimaryKeys(id, name)
	//studentModel.Find()

	//studentModel.UpdateByWhere(where, map[string]interface{}{
	//	"class": "vwwvwv",
	//	"Grade": "beebe",
	//	"from":  23,
	//})
	update := models.Student{
		Class: "",
		Grade: sql.NullString{String: "", Valid: true},
		From:  0,
	}
	studentModel.UpdateByWhere(where, update)
	//var list []models.Student
	//list = append(list, models.Student{Name: "vasv", Class: "vwqwv", From: 3})
	//list = append(list, models.Student{Name: "vasv", Class: "vwqwv", From: 4})
	//list = append(list, models.Student{Name: "vasv", Class: "vwqwv", From: 44})
	//studentModel.BatchInsert(list)
	return nil, nil
}

func (c *testService) Test3(id int) (string, error) {
	//if id%2 == 0 {
	//	time.Sleep(time.Second * 2)
	//} else if id%3 == 0 {
	//	time.Sleep(time.Second * 3)
	//} else if id%4 == 0 {
	//	time.Sleep(time.Second * 4)
	//} else if id%5 == 0 {
	//	time.Sleep(time.Second * 5)
	//} else {
	//time.Sleep(time.Second * 1)
	//}

	return strconv.Itoa(id), nil
}

func (c *testService) Test4(name string) ([]models.GroupBy, error) {
	a := &models.Student{}
	var b interface{}
	b = a
	err := gormdb.DefaultDB().Where("name = ?", name).Delete(b).Error
	return nil, err
}

func (c *testService) Test5(name string) ([]models.GroupBy, error) {
	err := gormdb.DefaultDB().Transaction(func(tx *gorm.DB) error {
		a := models.Student{}
		err := tx.Where("name = ?", name).First(&a).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if a.Id != 0 {
			itoa := strconv.Itoa(rand.Intn(50))
			err := tx.Table("student").Where("name = ?", name).
				Update("class", "csko"+itoa).Error
			if err != nil {
				return err
			}
		} else {
			a.Name = name
			tx.Create(&a)
		}
		return nil
	})
	return nil, err
}
