/*
Evo Package
Written by Spencer Liberto on 17 Sep 2015
*/

package main

import "fmt"
import "math/rand"
import "sort"
import "time"
import "strconv"

const logFile string = "/tmp/evo.go.log"
const numIndividualsPerColony int = 10
const numColonies int = 10
const numIndividuals int = numIndividualsPerColony * numColonies
const numGenes int = 30
const numGenerations int = 300
const childrenPerGeneration int = 2
const winnersPerGeneration int = childrenPerGeneration * 2 //winnersPerGeneration must be a multiple of numIndividuals

var r = rand.New(rand.NewSource(int64(time.Now().Unix())))

//Individual represents one Individual from one generation
type Individual struct {
  id int
  parentids [2]int
  fitnessScore float64
  genome [numGenes]float64
  tempRandomSorter int
}

//sort Individuals by fitnessScore
type ByFitnessScore []Individual
func (a ByFitnessScore) Len() int {return len(a)}
func (a ByFitnessScore) Swap(i, j int) {a[i], a[j] = a[j], a[i]}
func (a ByFitnessScore) Less(i, j int) bool {return a[i].fitnessScore < a[j].fitnessScore}

//sort Individuals by tempRandomSorter
type BySorter []Individual
func (a BySorter) Len() int {return len(a)}
func (a BySorter) Swap(i, j int) {a[i], a[j] = a[j], a[i]}
func (a BySorter) Less(i, j int) bool {return a[i].tempRandomSorter < a[j].tempRandomSorter}

func (i Individual) String() string {
  var s string
  s += strconv.Itoa(i.id) + " " + strconv.Itoa(i.parentids[0]) + " " + strconv.Itoa(i.parentids[1]) + " " + strconv.FormatFloat(i.fitnessScore, 'f', 6, 64) + " "
  for j:=0; j<numGenes; j++ {
    s+=strconv.FormatFloat(i.genome[j], 'f', 6, 64)+" "
  }
  return s
}

//Generates a random phenotype of all the Individuals
//This function is intended to be used to generate a first generation
func GenerateFirstGeneration(allIndividuals *[]Individual) {
  for i:=0; i<numIndividuals; i++ {
    (*allIndividuals)[i].id, (*allIndividuals)[i].parentids[0], (*allIndividuals)[i].parentids[1] = r.Int(), 0, 0
    for j:=0; j<numGenes; j++ {
      (*allIndividuals)[i].genome[j] = r.Float64()
    }
  }
}

//This function currently serves as an example game
//In this particular game, all genes are summed to produce a fitnessScore
func ColonyGame(colony []Individual, c chan bool) {
  for i:=0; i<numIndividualsPerColony; i++ {
    var g = colony[i].genome
    var sum float64
    for j:=0; j<numGenes; j++ {
      sum += g[j]
    }
    colony[i].fitnessScore = sum*100000
  }
  c <- true
}

func LaunchGeneration(newIndividuals []Individual) []Individual {
  c := make(chan bool, numColonies)
  for i:=0; i<numColonies; i++ {
    go ColonyGame(newIndividuals[numIndividualsPerColony*i:(numIndividualsPerColony*(i+1))], c)
  }
  for j:=0; j<numColonies; j++ {
    <-c
  }
  close(c)
  return newIndividuals
}

//Given an a list of Individuals, this function will choose the top fitnessScorers, randomly mate them, and return a slice of thier children
func ProduceNextGeneration (oldGen []Individual) (newGen []Individual) {
  //sort the old generation by fitnessScore. Winners will be at the top
  sort.Sort(sort.Reverse(ByFitnessScore(oldGen)))

  //Randomly sort winners by doing the following:
  //Copy winners to their own slice
  //assign winners random Individual.tempRandomSorter values and sort them
  var oldGenWinners []Individual
  oldGenWinners = make([]Individual, winnersPerGeneration)
  for i:=0; i<winnersPerGeneration; i++ {
    oldGenWinners[i] = oldGen[i]
    oldGenWinners[i].tempRandomSorter = r.Int()
  }
  sort.Sort(BySorter(oldGenWinners))

  //mate the winners to produce half as many children
  var uniqueChildren []Individual
  uniqueChildren = make([]Individual, childrenPerGeneration)
  for j:=0; j<childrenPerGeneration; j++ {
    uniqueChildren[j] = MateTwoIndividuals(oldGenWinners[j*2], oldGenWinners[j*2+1])
  }

  //make copies of children to fill up newGen
  //assign each one a unique random tempRandomSorter value for random sorting
  newGen = make([]Individual, numIndividuals)
  for k:=0; k<numIndividuals; k++ {
    newGen[k] = uniqueChildren[k%(winnersPerGeneration/2)]
    newGen[k].id = r.Int()
    newGen[k].tempRandomSorter = r.Int()
  }
  sort.Sort(BySorter(newGen))

  return newGen
}

//Child function of ProduceNextGeneration
//Given two Individuals, will produce a child
func MateTwoIndividuals (p1, p2 Individual) (child Individual) {
  child.parentids[0], child.parentids[1] = p1.id, p2.id

  var z int
  for geneNumber:=0; geneNumber<numGenes; geneNumber++ {
    z = r.Int()
    if z%200==0 || z%200==1 {
      child.genome[geneNumber] = r.Float64()
    } else if z%2==0 {
      child.genome[geneNumber] = p1.genome[geneNumber]
    } else {
      child.genome[geneNumber] = p2.genome[geneNumber]
    }
  }
  return child
}

func logToFileGeneration(allIndividuals []Individual, generationNumber int) {
  //to be implemented
  //for now you can enjoy your output on stdout
  fmt.Println("Generation",generationNumber)
  for i:=0; i<numIndividuals; i++ {
    fmt.Println(allIndividuals[i].String())
  }
}

func start(){
  var allIndividuals []Individual
  allIndividuals = make([]Individual, numIndividuals)
  GenerateFirstGeneration(&allIndividuals)
  for i:=0; i<numGenerations; i++ {
    allIndividuals = LaunchGeneration(allIndividuals)
    logToFileGeneration(allIndividuals, i)
    allIndividuals = ProduceNextGeneration(allIndividuals)
  }
}

func main() {
  start()
}
