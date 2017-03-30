package main

import "github.com/bgmerrell/simulago"

// simpy example:
/*
>>> import simpy
>>>
>>> def example(env):
...     value = yield env.timeout(1, value=42)
...     print('now=%d, value=%d' % (env.now, value))
>>>
>>> env = simpy.Environment()
>>> p = env.process(example(env))
>>> env.step()
now=1, value=42
*/

func example(env *simulago.Environment) {
}

func main() {
	env := simulago.NewEnvironment()
	_ = simulago.NewProcess(env, example)
	env.Step()
}
