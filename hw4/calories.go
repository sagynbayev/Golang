package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func BurnCalories(weight float64, height float64, age float64, activity string, gender string) float64 {
	var calories float64
	var active float64
	if activity == "a5" {
		active = 1.9
	} else if activity == "a4" {
		active = 1.725
	} else if activity == "a3" {
		active = 1.55
	} else if activity == "a2" {
		active = 1.375
	} else if activity == "a1" {
		active = 1.2
	}
	if gender == "male" || gender == "Male" {
		calories = (10*weight + 6.25*height - 5*age + 5) * active
	} else if gender == "female" || gender == "Female" {
		calories = (10*weight + 6.25*height - 5*age - 161) * active
	}
	return calories
}
func measurementDate() string {
	currentDate := time.Now()
	formattedDate := currentDate.Format("01-02-2006 15:04:05 Monday")
	return formattedDate
}

/*-----------------------------------------------------------------------------------------------------------------------*/
func MainMenu() {
	fmt.Println("Welcome to Calory Calculator :)")
	fmt.Println("Input 1 - to get action menu")
	fmt.Println("Input 0 - to close the program")
}
func ActionMenu() {
	fmt.Println("Input 1 - to calculate calories")
	fmt.Println("Input 2 - to record measurement date")
}
func main() {
	type infoJson struct {
		FirstName string
		LastName  string
		Calories  float64
		Date      string
	}
	type DateJson struct {
		FirstName string
		LastName  string
		Date      string
	}
	for {
		var choice1 int
		MainMenu()
		fmt.Scanln(&choice1)
		if choice1 == 1 {
			ActionMenu()
			var choice2 int
			fmt.Scanln(&choice2)
			if choice2 == 1 {
				var weight float64
				var height float64
				var age float64
				var activity string
				var gender string
				fmt.Println("Input your wieght")
				fmt.Scanln(&weight)
				fmt.Println("Input your height")
				fmt.Scanln(&height)
				fmt.Println("Input your age")
				fmt.Scanln(&age)
				fmt.Printf("Input your activity.\n There are 5 types of activity.\n Input a1 if minimal activity;\n Input a2 if weak activity;\n Input a3 if average activity;\n Input a4 if high activity;\n Input a5 if extra activity;\n")
				fmt.Scanln(&activity)
				fmt.Println("Input your gender.")
				fmt.Scanln(&gender)
				calories := BurnCalories(weight, height, age, activity, gender)
				fmt.Printf("We calculate number of calories that you burn everyday. Answer is %vcal \n \n", calories)
				var innerChoice int
				fmt.Println("Input 1 - to save this")
				fmt.Println("Input 0 - to go to the main menu")
				fmt.Scanln(&innerChoice)
				if innerChoice == 1 {
					currentDate := measurementDate()
					var firstName string
					var lastName string
					fmt.Println("Input your first name")
					fmt.Scanln(&firstName)
					fmt.Println("Input your last name")
					fmt.Scanln(&lastName)
					info := infoJson{
						FirstName: firstName,
						LastName:  lastName,
						Calories:  calories,
						Date:      currentDate,
					}
					b, err := json.Marshal(info)
					if err != nil {
						fmt.Println("error:", err)
					}
					os.Stdout.Write(b)
					fmt.Printf("\n\n")
				} else if innerChoice == 0 {
					fmt.Println("Okay, we going to menu without saving information about you")
				}
			} else if choice2 == 2 {
				currentDate := measurementDate()
				var firstName string
				var lastName string
				fmt.Println("Input your first name")
				fmt.Scanln(&firstName)
				fmt.Println("Input your last name")
				fmt.Scanln(&lastName)
				infoDate := DateJson{
					FirstName: firstName,
					LastName:  lastName,
					Date:      currentDate,
				}
				b, err := json.Marshal(infoDate)
				if err != nil {
					fmt.Println("error:", err)
				}
				os.Stdout.Write(b)
				fmt.Printf("\n\n")
			}
		} else if choice1 == 0 {
			break
		}
	}
}
