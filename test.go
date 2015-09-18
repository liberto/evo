package main

import "./evo"

//This function currently serves as an example game
//In this particular game, all genes are summed to produce a fitnessScore
func ColonyRound(colony []evo.Individual, c chan bool) {
  for i:=0; i<evo.NumIndividualsPerColony; i++ {
    var g = colony[i].Genome
    var sum float64
    for j:=0; j<evo.NumGenes; j++ {
      sum += g[j]
    }
    colony[i].FitnessScore = sum*100000
  }
  c <- true
}

func main(){
  evo.AssignColonyRound(ColonyRound)
  evo.StartGame()
}