package main

// type Cube struct {
// 	ID int
// }

// func (c *Cube) Render() {
// 	fmt.Printf("Renderuję kostkę o ID: %d\n", c.ID)
// }

// func main() {
// 	renderChan := make(chan *Cube)

// 	go func() {
// 		for {
// 			cube := <-renderChan

// 			if cube != nil {
// 				cube.Render()
// 			}
// 		}
// 	}()

// 	// temporary cube generation
// 	id := 1
// 	for {
// 		cube := &Cube{ID: id}
// 		id++

// 		renderChan <- cube

// 		time.Sleep(16 * time.Millisecond)
// 	}
// }
