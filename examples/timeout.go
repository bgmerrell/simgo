package main

import (
	"fmt"

	"github.com/bgmerrell/simgo"
	"github.com/bgmerrell/simgo/pcomm"
)

// simpy example:
/*
import simpy

def example(env):
    for i in range(2):
        event = simpy.events.Timeout(env, delay=10, value=42)
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
now=10, value=42
now=20, value=42
*/

func example(env *simgo.Environment, pc *pcomm.PCommunicator) {
	for i := 0; i < 2; i++ {
		to := simgo.NewTimeout(env, 10, 42)
		to.Schedule(env)
		pc.Send(to.Event)
		val, _ := to.Event.Value.Get()
		fmt.Printf("now=%d, value=%d\n", env.Now, val)
	}
}

func main() {
	env := simgo.NewEnvironment()
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, example))
	p.Init()
	/*
		_, err := env.Run(30)
			if err != nil {
				fmt.Print("Error: ", err)
			}
	*/
	env.Step()
	env.Step()
	env.Step()
}
