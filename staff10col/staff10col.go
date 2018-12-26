package staff10col

import (
	"database/sql"
	"github.com/badoux/checkmail"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/goFileImporter/file-importer/types"
	"github.com/yunabe/easycsv"
	"io"
	"log"
	"reflect"
	"regexp"
)

var RequireUnicode = []validation.Rule{
	validation.Required,
	validation.Match(regexp.MustCompile(`^[\d\p{L}\s’().'\-&",#_@+/]+$`)),
}

var LevelValidation = []validation.Rule{
	validation.Required,
	validation.Match(regexp.MustCompile(`^(teacher|superintendent|principal|super|admin|[2-4])$`)),
}

var RoleValidation = []validation.Rule{
	validation.Required,
	validation.Match(regexp.MustCompile(`inspect manager|evaluator|observer|student plan manager|staff`)),
}

// Staff - this is the struct for a staff
type Staff struct {
	FirstName    string     `name:"FirstName" json:"firstName" faker:"first_name"`
	LastName     string     `name:"LastName" json:"lastName" faker:"last_name"`
	Email        StaffEmail `name:"Email" json:"email" faker:"email"`
	Level        string     `name:"Level" json:"level"`
	Username     string     `name:"Username" json:"username"`
	Password     string     `name:"Password" json:"-"`
	SPN          string     `name:"SPN" json:"spn"`
	BuildingCode string     `name:"BuildingCode" json:"buildingCode"`
	BuildingName string     `name:"BuildingName" json:"buildingName"`
	Role         string     `name:"Role" json:"role"`
	action       string
}

type StaffEmail string

func (se StaffEmail) Valid() error {
	err := checkmail.ValidateFormat(string(se))
	if err != nil {
		return err
	}
	return nil
}

func (s Staff) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.FirstName, RequireUnicode...),
		validation.Field(&s.LastName, RequireUnicode...),
		validation.Field(&s.Username, RequireUnicode...),
		validation.Field(&s.Level, LevelValidation...),
		validation.Field(&s.Role, RoleValidation...),
	)
}

func GetDecoders() map[reflect.Type]interface{} {
	return map[reflect.Type]interface{}{
		reflect.TypeOf((StaffEmail)("")): func(s string) (StaffEmail, error) {
			return StaffEmail(s), nil
		},
	}
}

// StaffManager - this will house the configuration and the methods for working with a staff of many staff
type StaffManager struct {
	data           []types.Data
	header         bool
	reader         io.Reader
	erroredRecords []types.ErroredRecord
	dataStore      *sql.DB
}

func (s Staff) Valid() []error {
	// Go through validation process here
	var errs []error
	if err := s.Email.Valid(); err != nil {
		errs = append(errs, err)
	}
	if err := s.Validate(); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (sm *StaffManager) ValidateCollection() []types.ErroredRecord {
	for _, staff := range sm.data {
		if errs := staff.(Staff).Valid(); errs != nil {
			sm.erroredRecords = append(sm.erroredRecords, types.ErroredRecord{Err: errs, Data: staff.(Staff)})
		}
	}
	return sm.erroredRecords
}

// NewStaffManager - Constructor method for StaffManager
func NewStaffManager(db *sql.DB) *StaffManager {
	return &StaffManager{
		header:    true,
		dataStore: db,
	}
}

func NewEasyCSVReader(r io.Reader, decoders map[reflect.Type]interface{}) *easycsv.Reader {
	reader := easycsv.NewReader(r,
		easycsv.Option{
			TypeDecoders: decoders,
		},
	)
	return reader
}

func (sm *StaffManager) LoadDataFromReader(reader io.Reader) ([]types.Data, error) {
	var rows []types.Data
	r := NewEasyCSVReader(reader, GetDecoders())
	err := r.Loop(func(row Staff) error {
		rows = append(rows, row)
		return nil
	})
	if err != nil {
		//We could do something but better to let it trickle up
		return rows, err
	}
	sm.SetData(rows)
	return rows, err
}

func (sm *StaffManager) SetData(data []types.Data) {
	sm.data = data
}

// ShowData - return data structure
func (sm StaffManager) Data() []types.Data {
	return sm.data
}

func (sm StaffManager) ProcessData() []error {
	var errs []error
	db := sm.dataStore
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO users(username, email, first_name, last_name, level, user_key) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!
	for _, d := range sm.Data() {
		staff, _ := d.(Staff)
		_, err = stmt.Exec(staff.Username, staff.Email, staff.FirstName, staff.LastName, staff.Level, staff.SPN)
		if err != nil {
			errs = append(errs, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return errs
}
