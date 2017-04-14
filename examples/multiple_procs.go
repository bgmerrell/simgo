package main

import (
	"fmt"

	"github.com/bgmerrell/simgo"
)

// simpy example:
/*
import simpy

def my_proc(env):
    yield env.timeout(1)
    return 42

def other_proc(env):
    ret_val = yield env.process(my_proc(env))
    assert ret_val == 42

env = simpy.Environment()
gen = other_proc(env)
p = simpy.events.Process(env, gen)

env.run(p)
*/

func myProc(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	to := simgo.NewTimeout(env, 1, nil)
	to.Schedule(env)
	pc.Yield(to.Event)
	return 42
}

func otherProc(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, myProc))
	p.Init()
	pc.Yield(p.Event)
	fmt.Println("return value:", p.ReturnValue())
	return nil
}

func main() {
	env := simgo.NewEnvironment()
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, otherProc))
	p.Init()
	env.Run(nil)
}
