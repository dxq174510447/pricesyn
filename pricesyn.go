package main

import "fmt"

/**

 */
func main() {
	aa(12)
}

func aa(m int) (flowerror error){

	defer func() {
		if flowerror != nil {
			fmt.Printf("%v",flowerror)
		}else{
			fmt.Println("flowerror is nil")
		}
	}()

	if m < 10 {
		return nil
	}else{
		return fmt.Errorf("%d error",m)
	}


}
