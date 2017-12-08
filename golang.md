##### Arrays: 
* scores := []int{1,2,3,4,5} // size 5, capacity 5
* scores := make([]int, 10) // size 10, capacity 10
* scores := make([]int, 0, 10) // size 0, capacity 10
* scores = append(scores, 1) // may re-allocate
* c := cap(scores) // capacity
* for index, score := range scores {}
* copy(dest, source[:5])

##### Strings: 
* a := "the spice must flow"
* b := []byte(a)
* c := strinb(b)

###### Maps: 
* lookup := make(map[string]int)
* total := len(lookup)
* delete(lookup, "goku")
* lookup := make(map[string]int, 100)
* lookup := map[string]int{"goku":9001, "gohan":2044,}
* for key,value := range lookup {}

##### Lambdas: 
* type Add func(a int, b int) int
* func process (adder Add) int { return adder(1, 2) }
* process(func(a int, b int) int { return a + b})


##### Control: 
* initialized if : if err := process(); err != nil {}
* any: func add(a interface{}, b interface{} interface{} {}
* type-cast : return a.(int) + b.(int)
* type-switch: switch a.(type) { case int: case bool, string: default:}

##### Concurrency: 
* go function_name()
* go func() { println("this is lambda") }()
* c := make(chan int)
* c := make(chan int, 100)

```
for {
  select {
    case c <- rand.Int():
      // data successfully send
    default:
      // channel if full, data dropped
   }
}
```

```
for {
  select {
    case c <- rand.Int():
      // data sent
    case <-time.AFter(time.Millisecond * 100):
      // timed out
  }
}
```

##### Select: 
* the first available channel is chose 
* if multiple channels are available, one is randomly picked
* if no channle is available, the default case is executed
* if there's no default, block

##### Misc: 
* pass slices by value
* pass maps by reference
* pass structs by reference
* if the name starts with an uppercase char - it's visible outside the package
* if the nmae start with a lowercase char - it's not visigble outside the package


