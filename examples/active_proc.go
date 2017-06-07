package main

import (
	"fmt"

	"github.com/bgmerrell/simgo"
)

// simpy example:
/*
import simpy

def subfunc(env):
    print(env.active_process)  # will print "p1"

def my_proc(env):
    while True:
        print(env.active_process)  # will print "p1"
        subfunc(env)
        yield env.timeout(1)

env = simpy.Environment()
p1 = env.process(my_proc(env))
env.active_process  # None
env.step()
env.active_process  # None
*/

// Output:
/*
<Process(my_proc) object at 0x7f045064dac8>
<Process(my_proc) object at 0x7f045064dac8>
*/

func subfunc(env *simgo.Environment) {
	fmt.Printf("%p\n", env.ActiveProcess)
}

func myProc(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	for {
		fmt.Printf("%p\n", env.ActiveProcess)
		subfunc(env)
		to := simgo.NewTimeout(env, 1, nil)
		to.Schedule(env)
		pc.Yield(to.Event)
	}
	return nil
}

func main() {
	env := simgo.NewEnvironment()
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, myProc))
	p.Init()
	fmt.Printf("%p\n", env.ActiveProcess)
	env.Step()
	fmt.Printf("%p\n", env.ActiveProcess)
}
