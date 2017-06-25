package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Employee struct {
	EmployeeName string `json:"name"`
	Id           int    `json:"id"`
	ManagerId    int    `json:"managerid"`
}

func getEmployees() []Employee {
	//raw, err := ioutil.ReadFile("./employee-data.json")
	raw, err := ioutil.ReadFile("./employee-data-testing-with-errors.json")

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

func isEmployeeOK(emp Employee) bool {

	// check for errors first.  Will just report on them and carry on.  No need to crash the whole thing.
	switch {
	case emp.Id == 0:
		fmt.Println("\nERROR: Employee ID mssing: " + emp.toString())
		return false
	case emp.EmployeeName == "":
		fmt.Println("\nERROR: Employee name missing: " + emp.toString())
		return false
	}
	return true
}

func handler(w http.ResponseWriter, r *http.Request) {

	employees := getEmployees()
	for _, emp := range employees {

		if emp.ManagerId == 0 {
			if isEmployeeOK(emp) {
				fmt.Fprintf(w, "\n CEO: " + emp.EmployeeName)
			}
			var managers = getEmployeesForId(emp.Id)
			for _, mngs := range managers {
				if isEmployeeOK(emp) {
					fmt.Fprintf(w, "\n\t Managers: " + mngs.EmployeeName)
				}
				var minions = getEmployeesForId(mngs.Id)
				for _, m := range minions {
					if isEmployeeOK(m) {
						fmt.Fprintf(w, "\n\t\t Minions: " + m.EmployeeName)
					}
				}
			}
		}
	}
	fmt.Fprintf(w, "\n\n\n A copy of the input dataset can be found here: \n" + "https://storage.googleapis.com/nab-momenton-employee-datasets/employee-data.json\n\n")
}

func (employee Employee) toString() string {
	return toJson(employee)
}

func toJson(employee interface{}) string {
	bytes, err := json.Marshal(employee)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
