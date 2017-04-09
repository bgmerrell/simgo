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
    for i in range(3):
        event = simpy.events.Timeout(env, delay=10, value=40+1)
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
now=10, value=40
now=20, value=41
now=30, value=42
*/

func example(env *simgo.Environment, pc *pcomm.PCommunicator) {
	for i := 0; i < 3; i++ {
		to := simgo.NewTimeout(env, 10, 40+i)
		to.Schedule(env)
		pc.Send(to.Event)
		// BUG: Send() doesn't block here, so print statement below
		// can print env.Now value before it's updated.
		//time.Sleep(100 * time.Millisecond) // DELETE ME
		val, _ := to.Event.Value.Get()
		fmt.Printf("now=%d, value=%d\n", env.Now, val)
	}
}

func main() {
	env := simgo.NewEnvironment()
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, example))
	p.Init()
	/*
		_, err := env.Run(nil)
		if err != nil {
			fmt.Print("Error: ", err)
		}
	*/
	env.Step()
	env.Step()
	env.Step()
	env.Step()
}
