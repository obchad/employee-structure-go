package main

// Written by Chad OBrien
// I started learning GO a couple of weeks ago and saw this as an opportunity to write in GO rather than
// Java.
//  I have it hosted as a google app at https://nab-momenton-coding-challenge.appspot.com/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Employee object, loaded from a JSON file.
type Employee struct {
	EmployeeName string `json:"name"`
	Id           int    `json:"id"`
	ManagerId    int    `json:"managerid"`
}

// Get the employees from a JSON file.
func getEmployees() []Employee {
	raw, err := ioutil.ReadFile("./employee-data.json")
	//raw, err := ioutil.ReadFile("./employee-data-testing-with-errors.json")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var emp []Employee
	json.Unmarshal(raw, &emp)
	return emp
}

func getEmployeesForId(managerId int) []Employee {

	var minions []Employee

	if managerId > 0 {
		employees := getEmployees()
		for _, emp := range employees {

			if emp.ManagerId == managerId {
				minions = append(minions, emp)
			}
		}
	} else {
		// Console Output
		fmt.Println("WARN:  No manager id supplied, but will not fail here because that is not nice.  Plus it could be the ceo")
	}

	return minions
}

func isEmployeeOK(emp Employee, w http.ResponseWriter) bool {

	// Check for errors first.
	// Will just report on them and carry on.  No need to crash the whole thing.
	// I am outputting the console errors to the screen so you can see them easily even though it is ugly.
	// Otherwise I would ignore thos entries on the screen.
	switch {
	case emp.Id == 0:
		//fmt.Println("\nERROR: Employee ID mssing: " + emp.toString())
		fmt.Fprintf(w, "\n ----- Ugly ERROR message: Employee ID mssing: " + emp.toString())
		return false
	case emp.EmployeeName == "":
		//fmt.Println("\nERROR: Employee name missing: " + emp.toString())
		fmt.Fprintf(w, "\n ----- Ugly ERROR message: Employee name missing: " + emp.toString())
		return false
	}
	return true
}

// The HTTP url handler
func handler(w http.ResponseWriter, r *http.Request) {

	employees := getEmployees()
	var numberOfEmployees int = len(employees)
	fmt.Println("INFO: Number of employees in dataset: ", numberOfEmployees)

	// We want to report on the number or orphaned employees
	var numberOfDisplayedEmployeesCounter = 0

	for _, emp := range employees {

		if emp.ManagerId == 0 {
			if isEmployeeOK(emp, w) {
				fmt.Fprintf(w, "\n CEO: " + emp.EmployeeName)
				numberOfDisplayedEmployeesCounter++
			}
			var managers = getEmployeesForId(emp.Id)
			for _, mngs := range managers {
				if isEmployeeOK(emp, w) {
					fmt.Fprintf(w, "\n\t Managers: " + mngs.EmployeeName)
					numberOfDisplayedEmployeesCounter++
				}
				var minions = getEmployeesForId(mngs.Id)
				for _, m := range minions {
					if isEmployeeOK(m, w) {
						fmt.Fprintf(w, "\n\t\t Minions: " + m.EmployeeName)
						numberOfDisplayedEmployeesCounter++
					}
				}
				// If we needed to go deeper than 3 levels deep then I would start looking at parent
				// child tree implementation to handle infinite depth.  This takes more time
				// for me to do.  I have stuck to 3 levels here and to save time.
				// Ideally I would go back to the client and ask for more information on the spec
				// provided ie does it need to be flexible to handle more than 3 levels or is 3 ok.
			}
		}
	}
	// Console
	fmt.Println("INFO: Number of employees displayed: %s ", numberOfDisplayedEmployeesCounter)

	if numberOfEmployees > numberOfDisplayedEmployeesCounter {
		fmt.Fprintf(w,"\n\n You have orphaned employees in you dataset.  ")
	}
	fmt.Fprintf(w, "\n\n\n A copy of the input dataset can be found here: \n" + "https://storage.googleapis.com/nab-momenton-employee-datasets/employee-data.json\n\n")
	//fmt.Fprintf(w, "\n\n\n A copy of the input dataset can be found here: \n" + "https://storage.googleapis.com/nab-momenton-employee-datasets/employee-data-testing-with-errors.json\n\n")
}

// A tostring function for my error console outputs
func (employee Employee) toString() string {
	return toJson(employee)
}

// A toJson function for my error console outputs
func toJson(employee interface{}) string {
	bytes, err := json.Marshal(employee)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func main() {
	// if run locally see localhost:8080
	//  I have it hosted as a google app at https://nab-momenton-coding-challenge.appspot.com/
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
