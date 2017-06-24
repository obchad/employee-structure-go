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

func getEmployees() []Employee {
	raw, err := ioutil.ReadFile("./employee-data.json")

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

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

	employees := getEmployees()
	for _, emp := range employees {

		if emp.ManagerId == 0 {

			fmt.Fprintf(w, "\n Ceo: " + emp.EmployeeName)

			var managers = getEmployeesForId(emp.Id)

			for _, mngs := range managers {
				fmt.Fprintf(w, "\n\t Managers: " + mngs.EmployeeName)

				var minions = getEmployeesForId(mngs.Id)

				for _, m := range minions  {
					fmt.Fprintf(w, "\n\t\t Minions: " + m.EmployeeName)

				}

			}
			//fmt.Fprintf(w, "\n Managers: " + toJson(getEmployeesForId(emp.Id)))

		}

		// Console Output
		//fmt.Println(emp.toString())

		// Web output
		//fmt.Fprintf(w, emp.toString())
	}

	// Console Output
	//fmt.Println(toJson(employees))

	// Web Output
	//fmt.Fprintf(w, toJson(employees))

	//fmt.Fprintf(w, "\n" + toJson(getEmployeesForManagerId(400)))

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
