package main

import (
  "errors"
  "fmt"
  "strings"
  "runtime"
)

func main() {
  runtime.GOMAXPROCS(4)
  fields := []Field{
    Field{"West80",active, sheep, 350, 100, 80},
    Field{"Berket80",active, cow, 25, 12, 80},
    Field{"Demarest60",active, goat, 240, 80, 80},
    Field{"Lost40",inactive, "", 0, 0, 40},
  }

  farm := Farm{fields}

  ch := make(chan bool, 1)
  ch <- true
  for range ch {
    if input, err := getInput(); err == nil {
      fmt.Println("")
      switch input {
      case "1":
        animalCache := newAnimalCache()
        go farm.generateSalesOffspringReport(animalCache, ch)
      case "2":
        animalCache := newAnimalCache()
        go farm.generateOperatingReport(animalCache, ch)
      case "3":
        go farm.generateFarmReport(ch)
      }
    } else {
      fmt.Println(err.Error())
      close(ch)
    }
    <- ch
  }

}

func getInput() (option string, err error) {
  fmt.Println("1) Generate Sales Offspring Report")
  fmt.Println("2) Generate Operating Report")
  fmt.Println("3) Generate Farm Report")
  fmt.Println("Please choose an option: ")

  fmt.Scanln(&option)

  if option != "1" && option != "2"  && option != "3" {
    err = errors.New("Invalid option selected")
  }

  return
}


type Animal string
const (
  cow Animal = "Cow"
  sheep Animal = "Sheep"
  goat Animal = "Goat"
)

type AnimalInfo struct{
  name Animal
  perAcre float64
  offspringValue float64
}

type AnimalCache map[Animal]AnimalInfo

func newAnimalCache() AnimalCache {
  cowInfo := AnimalInfo{cow, 0.55, 500.}
  sheepInfo := AnimalInfo{sheep, 8., 100.}
  goatInfo := AnimalInfo{goat, 6., 200.}
  result := map[Animal]AnimalInfo{
    cow: cowInfo,
    sheep: sheepInfo,
    goat: goatInfo,
  }
  return result
}

type FieldStatus string

const (
  active FieldStatus = "Active"
  inactive FieldStatus = "Inactive"
)

type Field struct {
  name string
  status FieldStatus
  animal Animal
  adults uint
  offspring uint
  acreage uint
}

type Farm struct {
  fields []Field
}

func (f *Farm) generateSalesOffspringReport(animalCache AnimalCache, ch chan bool) {
  farmOffspringPotentialSales := 0.

  for _, field := range f.fields {
    fieldOffspringPotentialSales := 0.
    fieldTitle := fmt.Sprintf("Field Name %s:", field.name)
    fmt.Println(fieldTitle)
    fmt.Println(strings.Repeat("-", len(fieldTitle)))
    fieldOffspringPotentialSales = float64(field.offspring) * animalCache[field.animal].offspringValue
    farmOffspringPotentialSales += fieldOffspringPotentialSales

    fmt.Println("Animal: ", field.animal)
    fmt.Println("Offspring #: ", field.offspring)
    fmt.Println("Field Potential Sales: $", fieldOffspringPotentialSales)
    fmt.Println("")
  }
  farmTitle := "Farm Potential Sales"
  fmt.Println(farmTitle)
  fmt.Println(strings.Repeat("-", len(farmTitle)))
  fmt.Println("$", farmOffspringPotentialSales)
  ch <- true
  close(ch)
}


func (f *Farm) generateFarmReport(ch chan bool){
  for _, field := range f.fields {
    title := fmt.Sprintf("Field Name %s:", field.name)
    fmt.Println(title)
    fmt.Println(strings.Repeat("-", len(title)))
    fmt.Println("Acreage: ", field.acreage)
    if field.status == active {
      fmt.Println("Animal in field: ", field.animal)
      fmt.Println("Adult #: ", field.adults)
      fmt.Println("Offspring #: ", field.offspring)
      fmt.Println("")
    } else{
      fmt.Println("No livestock present")
    }
    fmt.Println("")
  }
  ch <- true
  close(ch)
}

func (f *Farm) generateOperatingReport(animalCache AnimalCache, ch chan bool) {
  farmAverageUtilization := 0.
  for _, field := range f.fields {
    fieldUtilization := 0.
    fieldTitle := fmt.Sprintf("Field Name: %s", field.name)
    fmt.Println(fieldTitle)
    fmt.Println(strings.Repeat("-", len(fieldTitle)))
    potential := 0.
    if field.status == active {
      potential = float64(field.acreage) * animalCache[field.animal].perAcre
      fieldUtilization = float64(field.adults)/potential*100.
    }
    farmAverageUtilization += fieldUtilization
    fmt.Printf("%-20s%.0f%s\n","Field Operating Utilization: ", fieldUtilization, "%")
    fmt.Printf("\n\n")
  }
  farmAverageUtilization = farmAverageUtilization/float64(len(f.fields))
  farmTitle := "Farm Average Operating Utilization"
  fmt.Println(farmTitle)
  fmt.Println(strings.Repeat("-", len(farmTitle)))
  fmt.Printf("Farm Operating Utilization: %.0f%s\n", farmAverageUtilization, "%")
  ch <- true
  close(ch)
}
