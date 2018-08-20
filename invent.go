// Project to demonstrate best place for inventory replisement
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// assume integers, units of inventory
// time units = 1 week

// randomly calculate the real consumption from estimation
func realConsumptionF(estimate int, variation int) int {

	rc := 0
	variationR := 0

	// Change seed for each execution

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	variationR = r1.Intn(variation)
	rc = estimate-variation/2+variationR

	// fmt.Println(" ******* Estimate ", estimate, " Variation ", variation/4, " VariationR ", variationR, " rc ", rc)

	return rc
}

// Forecast returns as percentage of installed base from consumption
func forecastReturnsF(fcastConsumptionM [104]int, percent float32) [104]int {

	var fcastReturnM [104]int

	for i := 0; len(fcastConsumptionM) > i; i++ {
		fcastReturnM[i] =  int(float32(fcastConsumptionM[i])*percent)
	}
	return fcastReturnM
}

// calculate and generate order, simple algorithm
func calcOrderSimple(inventory int, fcastConsumption [104]int, pendingOrder [104]int, orderDelay int, orderPeriod int, period int)  int {
	// inventoryPoint is the point when equipment is received
	var order, inventoryPoint int

	order = 0
	inventoryPoint = inventory

//	fmt.Println("Period ", period, " orderDelay ", orderDelay, " orderPeriod ", orderPeriod, " len ", len(fcastConsumption))

	// end at the data trail
	if ((period + orderDelay + orderPeriod) >= len(fcastConsumption)) {
		return order
	}

	// inventory point
	for i := 0; i < orderDelay; i++ {
		inventoryPoint =  inventoryPoint + pendingOrder[i+period] - fcastConsumption[i+period]
	}

	// add up forecast equipment until next reception
	for i := 0; i < orderPeriod; i++ {
		order =  order + pendingOrder[i+period+orderDelay] + fcastConsumption[i+period+orderDelay]
	}

	return order
}

func main() {
	// forecast consumption 24 months, per month. CPE number.
	fcastConsumption := [104]int{100, 200, 300, 350, 325, 450, 300, 250, 450, 475, 500, 250, //12
		100, 200, 300, 350, 325, 450, 300, 250, 450, 475, 500, 250, //24
		200, 250, 350, 450, 375, 550, 200, 270, 500, 575, 580, 350, //36
		100, 200, 300, 350, 325, 450, 300, 250, 450, 475, 500, 250, //48
		200, 250, 350, 450, 375, 550, 200, 270, 500, 575, 580, 350, //60
		100, 200, 300, 350, 325, 450, 300, 250, 450, 475, 500, 250, //72
		200, 250, 350, 450, 375, 550, 200, 270, 500, 575, 580, 350, //84
		200, 250, 350, 450, 375, 550, 200, 270, 500, 575, 580, 350, //96
		100, 200, 300, 350, 325, 450, 300, 250} //104

	// forecast returns 24 months, per month. CPE number.
	var fcastReturns, realConsumption [104]int

	// calculated consumption
	// var rcast [24]int

	//
	var totalConsumption, totalEstimated, totalReturned int
	var inventory, installBase  int
	var orderDelay int
	var orderPeriod int
	var cicloOrderDelay, cicloOrderPeriod int
	var pendingOrder [104]int

	// initial values
	inventory = 1200
	totalReturned = 0
	totalEstimated = 0
	totalConsumption = 0
	installBase = 5
	orderDelay = 6	// 1.5 months
	orderPeriod = 1 // 3 months
	cicloOrderDelay = 0
	cicloOrderPeriod = 0

	// calculate forecast returns
	fcastReturns = forecastReturnsF(fcastConsumption, 0.05)

//	for i := 0; i < len(fcastReturns); i++ {
//		fmt.Println("Month ", i, " return ", fcastReturns[i])
//	}

	for i := 0; 104 > i; i++ {
		realConsumption[i] = realConsumptionF(fcastConsumption[i], fcastConsumption[i]*2)
//		fmt.Println("Month ", i+1, " Estimation ", fcastConsumption[i], " Real ", realC, " Difference ", fcastConsumption[i]-realC)

		// update consumption, inventory, returned, installedBase
		totalEstimated = totalEstimated + fcastConsumption[i]
		totalConsumption = totalConsumption + realConsumption[i]
		inventory = inventory-realConsumption[i]+fcastReturns[i]
		totalReturned = totalReturned+fcastReturns[i]
		installBase = installBase+realConsumption[i]-fcastReturns[i]

		// evaluate reached conditions
		// equipment enter in warehouse
		if orderDelay*cicloOrderDelay == i {
			cicloOrderDelay++
			inventory = inventory+pendingOrder[i]
			fmt.Println("inventory ", inventory, " pendingOrder ", pendingOrder[i])
		}
		// calculate and order
		if orderPeriod*cicloOrderPeriod == i && i+orderDelay < len(fcastConsumption) {
			cicloOrderPeriod++
			pendingOrder[i+orderDelay] = calcOrderSimple(inventory, fcastConsumption, pendingOrder, orderDelay, orderPeriod, i)
//			fmt.Println("Period ", i+orderDelay, " Order ", pendingOrder[i+orderDelay], " i ", i)
		}
	}
	fmt.Println("Total Estimated ", totalEstimated, " Total Consumption ", totalConsumption, " Diference ", totalConsumption-totalEstimated)
	fmt.Println("Inventory ", inventory, " totalReturned ", totalReturned, " installBase ", installBase)
}