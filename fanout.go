package patterns

func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0, n)

	for i := 0; i < n; i++ {
		c := make(chan int)
		dests = append(dests, c)
		go func() {
			defer close(c)
			for n := range source {
				c <- n
			}
		}()
	}

	return dests
}
