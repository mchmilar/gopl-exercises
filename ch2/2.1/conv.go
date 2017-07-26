/*
»···Add types, constants, and functions for processing
»···temperatures in Kelvin, where zero Kelvin is -273.15C
»···and a difference of 1K has the same magnitude as 1C
*/

package tempconv

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// KToC converts a Kelvin temperature to a Celsius
func KToC(k Kelvin) Celsius { return Celsius(k - AbsoluteZeroC) }
