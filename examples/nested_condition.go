package main // simpy example:
import (
	"fmt"
	"log"

	"github.com/bgmerrell/simgo"
)

// Python example
/*
import simpy

def nested_condition(env):

    # Example nested `or` condition
    t1 = env.timeout(1, value='spam')
    t2 = env.timeout(2, value='eggs')
    t3 = env.timeout(3, value='coconut')
    results = yield (t1 & t2) | t3

    print("results: {}".format(results))
    assert results == {
        t1: 'spam',
        t2: 'eggs',
    }
    assert env.now == 2

    # Example nested `and` condition
    t1 = env.timeout(1, value='dog')
    t2 = env.timeout(9999, value='walrus')
    t3 = env.timeout(3, value='cat')
    results = yield (t1 | t2) & t3

    print("results: {}".format(results))
    assert results == {
        t1: 'dog',
        t3: 'cat',
    }
    assert env.now == 5


env = simpy.Environment()
proc = env.process(nested_condition(env))
env.run()

// Output
/*
results: <ConditionValue {<Timeout(1, value=spam) object at 0x7fe210468978>: 'spam', <Timeout(2, value=eggs) object at 0x7fe2104689b0>: 'eggs'}>
results: <ConditionValue {<Timeout(1, value=dog) object at 0x7fe210440e48>: 'dog', <Timeout(3, value=cat) object at 0x7fe210440eb8>: 'cat'}>

*/

func nestedCondition(env *simgo.Environment, pc *simgo.ProcComm) interface{} {

	// Example nested `or` condition
	t1 := simgo.NewTimeout(env, 1, "spam")
	t2 := simgo.NewTimeout(env, 2, "eggs")
	t3 := simgo.NewTimeout(env, 3, "coconut")
	t1.Schedule(env)
	t2.Schedule(env)
	t3.Schedule(env)

	events := []*simgo.Event{t1.Event, t2.Event}
	condition := simgo.AnyOf(
		env,
		[]*simgo.Event{simgo.AllOf(env, events).Event, t3.Event})
	r := pc.Yield(condition.Event).(simgo.ConditionValue)
	if len(r) != 2 {
		log.Fatalf("len(r) = %d, want: %d\n", len(r), 2)
	}
	val, err := r[0].(*simgo.EventValue).Get()
	fmt.Printf("val #1: %#v\n", val)
	if err != nil {
		log.Fatalf("err = %s, want: nil", err)
	}
	if val != "spam" {
		log.Fatalf("val = %s, want: spam", val)
	}
	val, err = r[1].(*simgo.EventValue).Get()
	fmt.Printf("val #2: %#v\n", val)
	if err != nil {
		log.Fatalf("err = %s, want: nil", err)
	}
	if val != "eggs" {
		log.Fatalf("val = %s, want: eggs", val)
	}

	// Example nested `and` condition
	fmt.Println("---")
	t1 = simgo.NewTimeout(env, 1, "cat")
	t2 = simgo.NewTimeout(env, 9999, "walrus")
	t3 = simgo.NewTimeout(env, 3, "dog")
	t1.Schedule(env)
	t2.Schedule(env)
	t3.Schedule(env)

	events = []*simgo.Event{t1.Event, t2.Event}
	condition = simgo.AllOf(
		env,
		[]*simgo.Event{simgo.AnyOf(env, events).Event, t3.Event})
	r = pc.Yield(condition.Event).(simgo.ConditionValue)
	if len(r) != 2 {
		log.Fatalf("len(r) = %d, want: %d\n", len(r), 2)
	}
	val, err = r[0].(*simgo.EventValue).Get()
	fmt.Printf("val #1: %#v\n", val)
	if err != nil {
		log.Fatalf("err = %s, want: nil", err)
	}
	if val != "cat" {
		log.Fatalf("val = %s, want: cat", val)
	}
	val, err = r[1].(*simgo.EventValue).Get()
	fmt.Printf("val #2: %#v\n", val)
	if err != nil {
		log.Fatalf("err = %s, want: nil", err)
	}
	if val != "dog" {
		log.Fatalf("val = %s, want: dog", val)
	}

	return nil
}

func main() {
	env := simgo.NewEnvironment()
	p := simgo.NewProcess(env, simgo.ProcWrapper(env, nestedCondition))
	p.Init()
	_, err := env.Run(nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
