// This package provides functionality to create .csv-files.
package invoicegen

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// CSV returns a string csv representation of the data supplied with genInputf
func CSV(genInput GeneratorInput) string {
	// Order should match order of row values below in `orderedRowValues`
	orderedColumnHeaders := []string{
		"Month",
		"ProjectID",
		"ContactPerson",
		"Amount",
		"Margin",
		"Total",
	}

	// Init csv writer
	content := new(bytes.Buffer)
	writer := csv.NewWriter(content)

	// write column headers
	writer.Write(orderedColumnHeaders)

	monthFormat := "2006-Jan"

	for _, row := range genInput.Entries {
		// Order should match order of columns above in `orderedColumnHeaders`
		orderedRowValues := []string{
			row.Month.Format(monthFormat),
			row.ProjectID,
			row.ContactPerson,
			fmt.Sprintf("%g", row.Amount),
			fmt.Sprintf("%g", row.Margin),
			fmt.Sprintf("%g", row.Total),
		}
		writer.Write(orderedRowValues)
	}

	writer.Flush()
	return content.String()
}

// DEPRECATED
// MarshalReflect tries to convert the passed value to csv using reflection.
// Returns an error if parsing using reflection is not possible.
// func MarshalReflect(genericSlice []interface{}) (string, error) {

// 	// Initialise csv writer
// 	content := new(bytes.Buffer)
// 	writer := csv.NewWriter(content)

// 	// perform type assurances on the first element of the slice
// 	t := reflect.TypeOf(genericSlice[0])
// 	switch t.Kind() {
// 	case reflect.Struct:

// 		fieldCount := t.NumField()

// 		orderedColumnHeaders := make([]string, fieldCount)
// 		orderedFieldNames := make([]string, t.NumField())

// 		for i := 0; i < fieldCount; i++ {
// 			field := t.Field(i)

// 			// ensure field implements the fmt.Stringer interface(i.e. , the type has a `func (self) String() string`)

// 			// Approaches inspired by
// 			// https://stackoverflow.com/questions/18570391/check-if-struct-implements-a-given-interface

// 			// Approach 1:
// 			_, ok := interface{}(field).(fmt.Stringer)
// 			if !ok && field.Type.String() != "string" {
// 				return "", errors.New("Provided field does not fulfill fmt.Stringer interface: " + field.Name)
// 			}

// 			// Approach 2:
// 			// stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
// 			// if !field.Type.Implements(stringerType) {
// 			// 	return "", errors.New("Provided field does not fulfill fmt.Stringer interface: " + field.Name)
// 			// }

// 			// ensure required field tags are present
// 			csvtag := field.Tag.Get("csv")
// 			if csvtag == "" {
// 				return "", errors.New("Non-empty struct tag 'csv' is required for marshaling struct. (describes desired csv header)")
// 			}
// 			order := field.Tag.Get("order")
// 			if csvtag == "" {
// 				return "", errors.New("Non-empty struct tag 'order' is required for marshaling struct. (describes desired csv column order)")
// 			}
// 			i, err := strconv.Atoi(order)
// 			if err != nil {
// 				return "", errors.New("Struct tag 'order' is required to be interpretable as an integer")
// 			}
// 			if i < 0 || i >= fieldCount {
// 				return "", errors.New("Values of struct tag 'order' have to be in the range [0, number of fields)")
// 			}

// 			// TODO: validate that i is in range(0, fieldCount) before using as index to avoid out of bounds panic

// 			orderedColumnHeaders[i] = csvtag
// 			orderedFieldNames[i] = field.Name
// 		}
// 		// type assurances have been made

// 		// print header row
// 		writer.Write(orderedColumnHeaders)

// 		// iterate over the slice and generate the csv presentation of each field
// 		for _, genericElement := range genericSlice {
// 			row := make([]string, fieldCount)
// 			for i, name := range orderedFieldNames {
// 				field := reflect.ValueOf(genericElement).FieldByName(name)
// 				row[i] = field.String()
// 			}
// 			err := writer.Write(row)
// 			if err != nil {
// 				log.Fatal("Could not write csv entry: ", err)
// 			}
// 		}
// 		writer.Flush()
// 	}
// 	return content.String(), nil
// }
