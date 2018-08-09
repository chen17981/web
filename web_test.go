
package main

import "testing"
import "math"


const Threshhold = 1e-9

func TestValidatePricelist(t *testing.T) {

	t.Run("Invalidate string with no number", func(t *testing.T){
	
		str := "CH1,CH2,CH3"
		_, _, got := validatePricelist(str)


		want := false
		if got != want {
			t.Errorf("got %t want %t given, %s", got, want, str)
		}
	})

	t.Run("Inalidate stirng with unmatched name and price pair", func(t *testing.T){
		
		str := "CH1,1.23,AP1,AP1"
		_, _, got := validatePricelist(str)
		
		want := false
		if got != want {
			t.Errorf("got %t want %t given, %s", got, want, str)
		}
		
	})
	
	t.Run("Inalidate stirng without product names", func(t *testing.T){
		
		str := "1.23,4,4,5"
		_, _, got := validatePricelist(str)
		
		want := false
		if got != want {
			t.Errorf("got %t want %t given, %s", got, want, str)
		}
		
	})

	t.Run("Validate string", func(t *testing.T){
	
		str := "CH1,1.23,AP1,1.43"
		s, n, got := validatePricelist(str)

		swant := [2]string{"CH1", "AP1"}
		nwant := [2]float64{1.23, 1.43}
		want  := true

		if got != want {	
			t.Errorf("got %t want %t given, %s", got, want, str)
		}
		
		for i, v := range swant {
			if s[i] != v {
				t.Errorf("got %s want %s given, %s", s[i], v, str)	
			}

			if n[i] != nwant[i] {
				t.Errorf("got %4.2f want %4.2f given, %s", n[i], nwant[i], str)
			}
		}
		
	})

}


func TestValidateShoppingItems(t *testing.T) {

	prices := map[string]float64 {
			"CH1" : 3.11,
			"AP1" : 6.00,
			"CF1" : 11.23,
			"MK1" : 4.75,
			"OM1" : 3.69,
			}
	
	t.Run("Invalid string with wrong product name", func(t *testing.T){
		
		str := "CB,CA,DA"
		_, got := validateShoppingItems(str, prices)

		want := false
		if got != want {
			t.Errorf("got %t want %t given, %s", got, want, str)
		}
	})

	t.Run("Invalid string with wrong seperator", func(t *testing.T){
		str := "CH1 AP1"
		
		_, got := validateShoppingItems(str, prices)

		want := false
		if got != want {
			t.Errorf("got %t want %t given, %s", got, want, str)
		}
	})	

	t.Run("Validate string", func(t *testing.T){
		str := "CH1,AP1"

		strs, got := validateShoppingItems(str, prices)

		want := true
		if got != want {
			t.Errorf("got %t want %t given, %s", got, want, str)
		}

		swant := [2]string{"CH1", "AP1"}
		for i,v := range swant {
			if strs[i] != v {
				t.Errorf("got %s want %s given, %s", strs[i], v, str)
			}
		}
	})	
}

func TestCalculate(t *testing.T) {

	prices := map[string]float64 {
			"CH1" : 3.11,
			"AP1" : 6.00,
			"CF1" : 11.23,
			"MK1" : 4.75,
			"OM1" : 3.69,
			}
	
	t.Run("Single item", func(t *testing.T){
		strs := []string{"CH1"}
		
		got := calculate(strs, prices, false, nil)

		want := float64(3.11)
                
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
	})

	t.Run("Two unrelated items", func(t *testing.T){
		strs := []string{"MK1", "AP1"}
		
		got := calculate(strs, prices, false, nil)
		
		want := float64(10.75)
                
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
		
	})

	t.Run("BOGO policy with two CF1s", func(t *testing.T){
		strs := []string{"CF1", "CF1"}
		
		got := calculate(strs, prices, false, nil)
		
		want := float64(11.23)
               
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
	})

	t.Run("BOGO policy with three CF1s", func(t *testing.T){
		strs := []string{"CF1", "CF1", "CF1"}
		
		got := calculate(strs, prices, false, nil)
		
		want := float64(22.46)
                
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
		
	})
	
	
	t.Run("BOGO policy with four CF1s", func(t *testing.T){
		strs := []string{"CF1", "CF1", "CF1", "CF1"}
		
		got := calculate(strs, prices, false, nil)
		
		want := float64(22.46)
                
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
		
	})

	t.Run("APPL policy with three apples and one unrelated item", func(t *testing.T){

		strs := []string{"AP1", "AP1", "CH1", "AP1"}
		got := calculate(strs, prices, false, nil)
		
		want := float64(16.61)
                if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
		
	})

	t.Run("APPL policy with three apples only", func(t *testing.T) {
		
		strs := []string{"AP1", "AP1", "AP1"}
		got := calculate(strs, prices, false, nil)
		
		want := float64(13.50)
                
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
		
	
	})

	t.Run("APPL policy with two apples only", func(t *testing.T) {
		
		strs := []string{"AP1", "AP1"}
		got := calculate(strs, prices, false, nil)
		
		want := float64(12.00)
                
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
		
	
	})
	
	t.Run("CHMK policy with four items", func(t *testing.T) {
		
		strs := []string{"CH1", "AP1", "CF1", "MK1"}
		got := calculate(strs, prices, false, nil)
		
		want := float64(20.34)
                
		if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }
			
	})


	 t.Run("APOM policy with four items", func(t *testing.T) {

                strs := []string{"AP1", "AP1", "OM1", "OM1"}
                got := calculate(strs, prices, false, nil)

                want := float64(13.38)

                if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }

        })

	
	t.Run("Sample test case", func(t *testing.T) {
		
		strs := []string{"CH1", "AP1", "AP1", "AP1", "MK1"}
		got := calculate(strs, prices, false, nil)
		
		want := float64(16.61)
                if math.Abs(got - want) > Threshhold {
                        t.Errorf("got %f want %f given, %s", got, want, strs)
                }	
	})	
}






