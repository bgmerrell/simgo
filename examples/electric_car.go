package main

import (
	"fmt"

	"github.com/bgmerrell/simgo"
)

// simpy example:
/*
class Car(object):
    def __init__(self, env):
        self.env = env
        # Start the run process everytime an instance is created.
        self.action = env.process(self.run())

    def run(self):
        while True:
            print('Start parking and charging at %d' % self.env.now)
            charge_duration = 5
            # We yield the process that process() returns
            # to wait for it to finish
            yield self.env.process(self.charge(charge_duration))

            # The charge process has finished and
            # we can start driving again.
            print('Start driving at %d' % self.env.now)
            trip_duration = 2
            yield self.env.timeout(trip_duration)

    def charge(self, duration):
        yield self.env.timeout(duration)

import simpy
env = simpy.Environment()
car = Car(env)
env.run(until=15)
*/

// Output
/*
Start parking and charging at 0
Start driving at 5
Start parking and charging at 7
Start driving at 12
Start parking and charging at 14
*/

const (
	charge_duration = 5
	trip_duration   = 2
)

type Car struct {
	env    *simgo.Environment
	action *simgo.Process
}

func NewCar(env *simgo.Environment) *Car {
	c := &Car{}
	c.env = env
	c.action = simgo.NewProcess(env, simgo.ProcWrapper(env, c.run))
	c.action.Init()
	return c
}

func (c *Car) run(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	for {
		fmt.Printf("Start parking and charging at %d\n", env.Now)

		to := simgo.NewTimeout(env, charge_duration, nil)
		to.Schedule(env)
		pc.Yield(to.Event)

		fmt.Printf("Start driving at %d\n", env.Now)
		to = simgo.NewTimeout(env, trip_duration, nil)
		to.Schedule(env)
		pc.Yield(to.Event)
	}
}

func (c *Car) charge(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	to := simgo.NewTimeout(env, charge_duration, nil)
	to.Schedule(env)
	pc.Yield(to.Event)
	return nil
}

func main() {
	env := simgo.NewEnvironment()
	car := NewCar(env)
	car.env.Run(15)
}
