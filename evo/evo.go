/*
Evo Package
Written by Spencer Liberto on 17 Sep 2015
*/

package evo

import "fmt"
import "math/rand"
import "sort"
import "time"
import "strconv"

var LogFile string = "/tmp/evo.go.log"
var NumIndividualsPerColony int = 10
var NumColonies int = 10
var NumIndividuals int = NumIndividualsPerColony * NumColonies
var NumGenes int = 30
var NumGenerations int = 10
var ChildrenPerGeneration int = 2
var WinnersPerGeneration int = ChildrenPerGeneration * 2 //WinnersPerGeneration must be a multiple of NumIndividuals

var r = rand.New(rand.NewSource(int64(time.Now().Unix())))
type ColonyRound func([]Individual, chan bool) 
var userDefinedColonyRound ColonyRound

//Individual represents one Individual from one generation
type Individual struct {
  Id int
  Parentids [2]int
  FitnessScore float64
  Genome [NumGenes]float64
  tempRandomSorter int
}

//sort Individuals by fitnessScore
type ByFitnessScore []Individual
func (a ByFitnessScore) Len() int {return len(a)}
func (a ByFitnessScore) Swap(i, j int) {a[i], a[j] = a[j], a[i]}
func (a ByFitnessScore) Less(i, j int) bool {return a[i].FitnessScore < a[j].FitnessScore}

//sort Individuals by tempRandomSorter
type BySorter []Individual
func (a BySorter) Len() int {return len(a)}
func (a BySorter) Swap(i, j int) {a[i], a[j] = a[j], a[i]}
func (a BySorter) Less(i, j int) bool {return a[i].tempRandomSorter < a[j].tempRandomSorter}

func (i Individual) String() string {
  var s string
  s += strconv.Itoa(i.Id) + " " + strconv.Itoa(i.Parentids[0]) + " " + strconv.Itoa(i.Parentids[1]) + " " + strconv.FormatFloat(i.FitnessScore, 'f', 6, 64) + " "
  for j:=0; j<NumGenes; j++ {
    s+=strconv.FormatFloat(i.Genome[j], 'f', 6, 64)+" "
  }
  return s
}

//Generates a random phenotype of all the Individuals
//This function is intended to be used to generate a first generation
func generateFirstGeneration(allIndividuals *[]Individual) {
  for i:=0; i<NumIndividuals; i++ {
    (*allIndividuals)[i].Id, (*allIndividuals)[i].Parentids[0], (*allIndividuals)[i].Parentids[1] = r.Int(), 0, 0
    for j:=0; j<NumGenes; j++ {
      (*allIndividuals)[i].Genome[j] = r.Float64()
    }
  }
}


func launchGeneration(newIndividuals []Individual) []Individual {
  c := make(chan bool, NumColonies)
  for i:=0; i<NumColonies; i++ {
    go userDefinedColonyRound(newIndividuals[NumIndividualsPerColony*i:(NumIndividualsPerColony*(i+1))], c)
  }
  for j:=0; j<NumColonies; j++ {
    <-c
  }
  close(c)
  return newIndividuals
}

//Given an a list of Individuals, this function will choose the top fitnessScorers, randomly mate them, and return a slice of thier children
func produceNextGeneration (oldGen []Individual) (newGen []Individual) {
  //sort the old generation by fitnessScore. Winners will be at the top
  sort.Sort(sort.Reverse(ByFitnessScore(oldGen)))

  //Randomly sort winners by doing the following:
  //Copy winners to their own slice
  //assign winners random Individual.tempRandomSorter values and sort them
  var oldGenWinners []Individual
  oldGenWinners = make([]Individual, WinnersPerGeneration)
  for i:=0; i<WinnersPerGeneration; i++ {
    oldGenWinners[i] = oldGen[i]
    oldGenWinners[i].tempRandomSorter = r.Int()
  }
  sort.Sort(BySorter(oldGenWinners))

  //mate the winners to produce half as many children
  var uniqueChildren []Individual
  uniqueChildren = make([]Individual, ChildrenPerGeneration)
  for j:=0; j<ChildrenPerGeneration; j++ {
    uniqueChildren[j] = mateTwoIndividuals(oldGenWinners[j*2], oldGenWinners[j*2+1])
  }

  //make copies of children to fill up newGen
  //assign each one a unique random tempRandomSorter value for random sorting
  newGen = make([]Individual, NumIndividuals)
  for k:=0; k<NumIndividuals; k++ {
    newGen[k] = uniqueChildren[k%(ChildrenPerGeneration)]
    newGen[k].Id = r.Int()
    newGen[k].tempRandomSorter = r.Int()
  }
  sort.Sort(BySorter(newGen))

  return newGen
}

//Child function of ProduceNextGeneration
//Given two Individuals, will produce a child
func mateTwoIndividuals (p1, p2 Individual) (child Individual) {
  child.Parentids[0], child.Parentids[1] = p1.Id, p2.Id

  var z int
  for geneNumber:=0; geneNumber<NumGenes; geneNumber++ {
    z = r.Int()
    if z%200==0 || z%200==1 {
      child.Genome[geneNumber] = r.Float64()
    } else if z%2==0 {
      child.Genome[geneNumber] = p1.Genome[geneNumber]
    } else {
      child.Genome[geneNumber] = p2.Genome[geneNumber]
    }
  }
  return child
}

func logToFileGeneration(allIndividuals []Individual, generationNumber int) {
  //to be implemented
  //for now you can enjoy your output on stdout
  fmt.Println("Generation",generationNumber)
  for i:=0; i<NumIndividuals; i++ {
    fmt.Println(allIndividuals[i].String())
  }
}

func StartGame(){
  var allIndividuals []Individual
  allIndividuals = make([]Individual, NumIndividuals)
  generateFirstGeneration(&allIndividuals)
  for i:=0; i<NumGenerations; i++ {
    allIndividuals = launchGeneration(allIndividuals)
    logToFileGeneration(allIndividuals, i)
    allIndividuals = produceNextGeneration(allIndividuals)
  }
}

func Hello(){
  fmt.Println("hello")
}

func AssignColonyRound(cr ColonyRound) {
  userDefinedColonyRound = cr
  fmt.Println("yes")
}

func main(){
  fmt.Println("nope")
}