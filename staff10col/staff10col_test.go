package staff10col

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/bxcodec/faker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goFileImporter/file-importer/types"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type StaffTestSuite struct {
	suite.Suite
	staff Staff
}

type StaffManagerSuite struct {
	suite.Suite
	staffManager      types.Manager
	staffManagerFaker types.Manager
}

type StaffManagerSuiteWithErrs struct {
	suite.Suite
	staffManager      types.Manager
	staffManagerFaker types.Manager
}

func (suite *StaffTestSuite) SetupTest() {
	err := faker.FakeData(&suite.staff)
	suite.staff.Level = "admin"
	suite.staff.Role = "staff"
	if err != nil {
		panic(err)
	}
}

func (suite *StaffTestSuite) TestStaffStructValid() {
	suite.Empty(suite.staff.Valid())
}

func TestStaffSuite(t *testing.T) {
	suite.Run(t, new(StaffTestSuite))
}

func (suite *StaffManagerSuite) SetupTest() {
	var staffManager types.Manager
	db, err := sql.Open("mysql", "fake:fake@/fake")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	staffManager = NewStaffManager(db)
	_, staffManagerOk := staffManager.(types.Manager)
	suite.True(staffManagerOk)
	suite.staffManager = staffManager

	var staff Staff
	for i := 0; i < 10; i++ {
		err := faker.FakeData(&staff)
		staff.Level = "admin"
		staff.Role = "evaluator"
		if err != nil {
			panic(err)
		}
		suite.staffManager.SetData(append(suite.staffManager.Data(), staff))
	}
}

func (suite *StaffManagerSuiteWithErrs) SetupTest() {
	var staffManager types.Manager
	db, err := sql.Open("mysql", "fake:fake@/fake")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err != nil {
		// This is probably bad
	}
	staffManager = NewStaffManager(db)
	_, staffManagerOk := staffManager.(types.Manager)
	suite.True(staffManagerOk)
	suite.staffManager = staffManager

	var staff Staff
	var st []types.Data
	for i := 0; i < 10; i++ {
		err := faker.FakeData(&staff)
		staff.Email = "bob.bob.com"
		if err != nil {
			panic(err)
		}
		st = append(st, staff)
	}

	suite.staffManager.SetData(st)
}

func (s *StaffManagerSuite) TestLoadDataFromPath() {
	var data []types.Data
	var err error

	file, err := os.Open("../testdata/staffManager.csv")
	data, err = s.staffManager.LoadDataFromReader(bufio.NewReader(file))
	if s.Nil(err) {
		s.Equal((data[1]).(Staff).FirstName, "John")
	}
}

func (suite *StaffManagerSuite) TestStaffCollectionValid() {
	errs := (suite.staffManager).ValidateCollection()
	for _, err := range errs {
		fmt.Println(err.Err)
	}
	suite.Empty(errs)
}

func (suite *StaffManagerSuiteWithErrs) TestStaffCollectionNotValid() {

	errs := (suite.staffManager).ValidateCollection()
	suite.NotEmpty(errs)
}

func (s *StaffManagerSuiteWithErrs) TestLoadDataFromPath() {
	var err error
	file, _ := os.Open("failing/file/that/will/fail/do/not/put/a/file/here.csv")
	_, err = s.staffManager.LoadDataFromReader(bufio.NewReader(file))
	s.NotNil(err)
}

func (suite *StaffManagerSuite) TestShowData() {
	suite.NotEmpty(suite.staffManager.Data())
}

func TestStaffManagerSuite(t *testing.T) {
	suite.Run(t, new(StaffManagerSuite))
	suite.Run(t, new(StaffManagerSuiteWithErrs))
}
