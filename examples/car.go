package main

import (
	"fmt"

	"github.com/bgmerrell/simgo"
)

// simpy example:
/*
def car(env):
    while True:
        print('Start parking at %d' % env.now)
        parking_duration = 5
        yield env.timeout(parking_duration)

        print('Start driving at %d' % env.now)
        trip_duration = 2
        yield env.timeout(trip_duration)

import simpy
env = simpy.Environment()
env.process(car(env))
env.run(until=15)
*/

// Output:
/*
Start parking at 0
Start driving at 5
Start parking at 7
Start driving at 12
Start parking at 14
*/

func car(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	const (
		parking_duration = 5
		trip_duration    = 2
	)
	for {
		fmt.Printf("Start parking at %d\n", env.Now)

		to := simgo.NewTimeout(env, parking_duration, nil)
		to.Schedule(env)
		pc.Yield(to.Event)

		fmt.Printf("Start driving at %d\n", env.Now)
		to = simgo.NewTimeout(env, trip_duration, nil)
		to.Schedule(env)
		pc.Yield(to.Event)
	}
}

func main() {
	env := simgo.NewEnvironment()
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, car))
	p.Init()
	env.Run(15)
}
