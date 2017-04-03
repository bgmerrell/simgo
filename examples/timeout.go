package main

import (
	"github.com/bgmerrell/simulago"
	"github.com/bgmerrell/simulago/pcomm"
)

// simpy example:
/*
import simpy

def example(env):
    for i in range(2):
        event = simpy.events.Timeout(env, delay=1, value=42)
        value = yield event
        print('now=%d, value=%d' % (env.now, value))

env = simpy.Environment()
example_gen = example(env)
_ = simpy.events.Process(env, example_gen)

env.step()
env.step()
env.step()
*/

// Output:
/*
now=1, value=42
now=2, value=42
*/

func example(env *simulago.Environment, pc *pcomm.PCommunicator) {
	for i := 0; i < 2; i++ {
		to := simulago.NewTimeout(env, 10)
		to.Schedule(env)
		pc.Send(to.Event)
	}
}

func main() {
	env := simulago.NewEnvironment()
	pc := simulago.ProcWrapper(env, example)
	p := simulago.NewProcess(env, pc)
	p.Init()
	env.Step()
	env.Step()
	env.Step()
}
