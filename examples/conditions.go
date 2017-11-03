package main // simpy example:
import (
	"fmt"
	"log"

	"github.com/bgmerrell/simgo"
)

/*
import simpy

def test_condition(env):
    t1, t2 = env.timeout(1, value='spam'), env.timeout(2, value='eggs')
    ret = yield t1 | t2
    assert ret == {t1: 'spam'}
    print("t1 finished")

    t1, t2 = env.timeout(1, value='spam'), env.timeout(2, value='eggs')
    ret = yield t1 & t2
    assert ret == {t1: 'spam', t2: 'eggs'}
    print("t1 and t2 finished")

proc = env.process(test_condition(env))
env.run()
*/

// Output
/*
t1 finished
t1 and t2 finished
*/

func testCondition(env *simgo.Environment, pc *simgo.ProcComm) interface{} {

	// AnyOf() example
	t1 := simgo.NewTimeout(env, 1, "spam")
	t2 := simgo.NewTimeout(env, 2, "eggs")
	t1.Schedule(env)
	t2.Schedule(env)

	events := []*simgo.Event{t1.Event, t2.Event}
	condition := simgo.AnyOf(env, events)
	r := pc.Yield(condition.Event).(simgo.ConditionValue)
	if len(r) != 1 {
		log.Fatalf("len(r) = %d, want: %d\n", len(r), 1)
	}
	val, err := r[0].(*simgo.EventValue).Get()
	fmt.Printf("AnyOf val #1: %#v\n", val)
	if err != nil {
		log.Fatalf("err = %s, want: nil", err)
	}
	if val != "spam" {
		log.Fatalf("val = %s, want: spam", val)
	}

	// AllOf() example
	t1 = simgo.NewTimeout(env, 1, "spam")
	t2 = simgo.NewTimeout(env, 2, "eggs")
	t1.Schedule(env)
	t2.Schedule(env)

	events = []*simgo.Event{t1.Event, t2.Event}
	condition = simgo.AllOf(env, events)
	r = pc.Yield(condition.Event).(simgo.ConditionValue)
	if len(r) != 2 {
		log.Fatalf("len(r) = %d, want: %d\n", len(r), 2)
	}
	val, err = r[0].(*simgo.EventValue).Get()
	fmt.Printf("AllOf() val #1: %#v\n", val)
	if err != nil {
		log.Fatalf("err = %s, want: nil", err)
	}
	if val != "spam" {
		log.Fatalf("val = %s, want: spam", val)
	}
	val, err = r[1].(*simgo.EventValue).Get()
	fmt.Printf("AllOfval val #2: %#v\n", val)
	if err != nil {
		log.Fatalf("err = %s, want: nil", err)
	}
	if val != "eggs" {
		log.Fatalf("val = %s, want: eggs", val)
	}

	return nil
}

func main() {
	env := simgo.NewEnvironment()
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, testCondition))
	p.Init()
	_, err := env.Run(nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
