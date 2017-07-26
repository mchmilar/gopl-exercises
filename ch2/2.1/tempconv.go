/*
	Add types, constants, and functions for processing
	temperatures in Kelvin, where zero Kelvin is -273.15C
	and a difference of 1K has the same magnitude as 1C
*/
package tempconv

import "fmt"

type Celsius float64
type Fahrenheit float64
type Kelvie float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }
func (k Kelvin) String() string     { return fmt.Sprintf("%gK", k) }
