package main

import (
	"fmt"

	"github.com/bgmerrell/simgo"
)

// simpy example:
/*
import simpy

class School:
    def __init__(self, env):
        self.env = env
        self.class_ends = env.event()
        self.pupil_procs = [env.process(self.pupil()) for i in range(3)]
        self.bell_proc = env.process(self.bell())

    def bell(self):
        for i in range(2):
            yield self.env.timeout(45)
            self.class_ends.succeed()
            self.class_ends = self.env.event()
            print()

    def pupil(self):
        for i in range(2):
            print(' \o/', end='')
            yield self.class_ends

env = simpy.Environment()
school = School(env)
env.run()
*/

// Output:
/*
 \o/ \o/ \o/
 \o/ \o/ \o/
*/

type School struct {
	env        *simgo.Environment
	classEnds  *simgo.Event
	pupilProcs []*simgo.Process
	bellProc   *simgo.Process
}

func (s *School) bell(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	for i := 0; i < 2; i++ {
		to := simgo.NewTimeout(env, 45, nil)
		to.Schedule(env)
		pc.Yield(to.Event)
		s.classEnds.Succeed(nil)
		s.classEnds = simgo.NewEvent(env)
		fmt.Println("")
	}
	return nil
}

func (s *School) pupil(env *simgo.Environment, pc *simgo.ProcComm) interface{} {
	for i := 0; i < 2; i++ {
		fmt.Printf(` \o/`)
		pc.Yield(s.classEnds)
	}
	return nil
}

func NewSchool(env *simgo.Environment) *School {
	s := &School{
		env:       env,
		classEnds: simgo.NewEvent(env),
	}
	s.bellProc = simgo.NewProcess(env, simgo.ProcWrapper(env, s.bell))
	s.bellProc.Init()
	s.pupilProcs = make([]*simgo.Process, 0, 3)
	for i := 0; i < 3; i++ {
		pupilProc := simgo.NewProcess(env, simgo.ProcWrapper(env, s.pupil))
		pupilProc.Init()
		s.pupilProcs = append(s.pupilProcs, pupilProc)
	}
	return s
}

func main() {
	env := simgo.NewEnvironment()
	school := NewSchool(env)
	school.env.Run(nil)
}
