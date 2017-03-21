package main

import (
  "math"
  "math/rand"
  "os"
  "strconv"
  "time"
  log "github.com/Sirupsen/logrus"
  gcd "github.com/gilgamez/algorithms/algorithms/maths/stein"
)

func generateRandomNumbers(max int) (chan int) {
  c := make(chan int, 1000)
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  go func() {
    for {
      value := 1 + r.Intn(max)
      c <- value
    }
  }()
  return c
}

type PiStimate struct {
  coprime int
  total   int
}

func (ps *PiStimate) Coprime() int {
  return ps.coprime
}

func (ps *PiStimate) Cofactor() int {
  return ps.total - ps.coprime
}

func (ps *PiStimate) Total() int {
  return ps.total
}

func (ps *PiStimate) Pi() float64 {
  return math.Sqrt(float64(6) / (float64(ps.coprime) / float64(ps.total)))
}

func (ps *PiStimate) Add(other PiStimate) {
  ps.total += other.total
  ps.coprime += other.coprime
}

var coprime = PiStimate{coprime: 1, total: 1 }
var cofactor = PiStimate{coprime: 0, total: 1 }

func decide(a int, b int) *PiStimate {
  if gcd.Iter(a, b) == 1 {
    return &coprime
  } else {
    return &cofactor
  }
}

func getArg(index int, defaultValue int) int {
  result := defaultValue
  if len(os.Args) > index {
    argValue, err := strconv.Atoi(os.Args[index])
    if err == nil {
      result = argValue
    }
  }
  return result
}

func main() {
  samples := getArg(1, 100000)
  upper := int(getArg(2, math.MaxInt64))
  a := generateRandomNumbers(upper)
  b := generateRandomNumbers(upper)

  pi := make(chan *PiStimate, 1000)

  go func() {
    for i := 0; i < samples; i++ {
      pi <- decide(<-a, <-b)
    }
    close(pi)
  }()

  result := PiStimate{coprime: 0, total: 0}
  log.WithFields(log.Fields{"samples": samples, "upper": upper}).Info("parameters")

  for v := range pi {
    result.Add(*v)
  }

  log.WithFields(log.Fields{"pi": result.Pi(), "total": result.Total()}).Info("result")
}
