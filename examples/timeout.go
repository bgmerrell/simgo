package main

import "github.com/bgmerrell/simulago"

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

func example(env *simulago.Environment) chan struct{} {
	ch := make(chan struct{})
	go func() {
		for i := 0; i < 2; i++ {
			<-ch
			to := simulago.NewTimeout(env, 10, 42)
			to.Schedule(env)
		}
	}()
	return ch
}

func main() {
	env := simulago.NewEnvironment()
	ch := example(env)
	_ = simulago.NewProcess(env, ch)
	// TODO: This is broken. Need to decide on how to best communicate
	// with process.
	ch <- struct{}{}
	env.Step()
}
