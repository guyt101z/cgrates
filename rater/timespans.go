/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2013 ITsysCOM

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package rater

import (
	"fmt"
	"math"
	"time"
)

/*
A unit in which a call will be split that has a specific price related interval attached to it.
*/
type TimeSpan struct {
	TimeStart, TimeEnd time.Time
	Cost               float64
	ActivationPeriod   *ActivationPeriod
	Interval           *Interval
	MinuteInfo         *MinuteInfo
}

// Holds the bonus minute information related to a specified timespan
type MinuteInfo struct {
	DestinationId string
	Quantity      float64
	Price         float64
}

/*
Returns the duration of the timespan
*/
func (ts *TimeSpan) GetDuration() time.Duration {
	return ts.TimeEnd.Sub(ts.TimeStart)
}

// Returns the cost of the timespan according to the relevant cost interval.
// It also sets the Cost field of this timespan (used for refound on session
// manager debit loop where the cost cannot be recalculated)
func (ts *TimeSpan) getCost(cd *CallDescriptor) (cost float64) {
	if ts.MinuteInfo != nil {
		return ts.GetDuration().Seconds() * ts.MinuteInfo.Price
	}
	if ts.Interval == nil {
		return 0
	}
	duration := ts.GetDuration().Seconds()
	i := ts.Interval
	if i.RateIncrements == 0 {
		i.RateIncrements = 1
	}
	if i.PricedUnits != 0 {
		cost = math.Ceil(duration/i.RateIncrements) * i.RateIncrements * (i.Price / i.PricedUnits)
	} else {
		cost = math.Ceil(duration/i.RateIncrements) * i.RateIncrements * i.Price
	}
	// if userBalance, err := cd.getUserBalance(); err == nil && userBalance != nil {
	// 	userBalance.mux.RLock()
	// 	if percentageDiscount, err := userBalance.getVolumeDiscount(cd.Destination, INBOUND); err == nil && percentageDiscount > 0 {
	// 		cost *= (100 - percentageDiscount) / 100
	// 	}
	// 	userBalance.mux.RUnlock()
	// }
	ts.Cost = cost
	return
}

/*
Returns true if the given time is inside timespan range.
*/
func (ts *TimeSpan) Contains(t time.Time) bool {
	return t.After(ts.TimeStart) && t.Before(ts.TimeEnd)
}

/*
Will set the interval as spans's interval if new Weight is greater then span's interval Weight
or if the Weights are equal and new price is lower then spans's interval price
*/
func (ts *TimeSpan) SetInterval(i *Interval) {
	if ts.Interval == nil || ts.Interval.Weight < i.Weight {
		ts.Interval = i
	}
	if ts.Interval.Weight == i.Weight && i.Price < ts.Interval.Price {
		ts.Interval = i
	}
}

/*
Splits the given timespan according to how it relates to the interval.
It will modify the endtime of the received timespan and it will return
a new timespan starting from the end of the received one.
The interval will attach itself to the timespan that overlaps the interval.
*/
func (ts *TimeSpan) SplitByInterval(i *Interval) (nts *TimeSpan) {
	//Logger.Debug("here: ", ts, " +++ ", i)
	// if the span is not in interval return nil	
	if !(i.Contains(ts.TimeStart) || i.Contains(ts.TimeEnd)) {
		//Logger.Debug("Not in interval")
		return
	}
	// if the span is enclosed in the interval try to set as new interval and return nil
	if i.Contains(ts.TimeStart) && i.Contains(ts.TimeEnd) {
		//Logger.Debug("All in interval")
		ts.SetInterval(i)
		return
	}
	// if only the start time is in the interval split the interval
	if i.Contains(ts.TimeStart) {
		//Logger.Debug("Start in interval")
		splitTime := i.getRightMargin(ts.TimeStart)
		ts.SetInterval(i)
		if splitTime == ts.TimeStart {
			return
		}
		nts = &TimeSpan{TimeStart: splitTime, TimeEnd: ts.TimeEnd}
		ts.TimeEnd = splitTime

		return
	}
	// if only the end time is in the interval split the interval
	if i.Contains(ts.TimeEnd) {
		//Logger.Debug("End in interval")
		splitTime := i.getLeftMargin(ts.TimeEnd)
		if splitTime == ts.TimeEnd {
			return
		}
		nts = &TimeSpan{TimeStart: splitTime, TimeEnd: ts.TimeEnd}
		ts.TimeEnd = splitTime

		nts.SetInterval(i)
		return
	}
	return
}

/*
Splits the given timespan on activation period's activation time.
*/
func (ts *TimeSpan) SplitByActivationPeriod(ap *ActivationPeriod) (newTs *TimeSpan) {
	if !ts.Contains(ap.ActivationTime) {
		return nil
	}
	newTs = &TimeSpan{TimeStart: ap.ActivationTime, TimeEnd: ts.TimeEnd, ActivationPeriod: ap}
	ts.TimeEnd = ap.ActivationTime
	return
}

/*
Splits the given timespan on minute bucket's duration.
*/
func (ts *TimeSpan) SplitByMinuteBucket(mb *MinuteBucket) (newTs *TimeSpan) {
	s := ts.GetDuration().Seconds()
	ts.MinuteInfo = &MinuteInfo{mb.DestinationId, s, mb.Price}
	if s <= mb.Seconds {
		mb.Seconds -= s
		return nil
	}
	secDuration, _ := time.ParseDuration(fmt.Sprintf("%vs", mb.Seconds))

	newTimeEnd := ts.TimeStart.Add(secDuration)
	newTs = &TimeSpan{TimeStart: newTimeEnd, TimeEnd: ts.TimeEnd}
	ts.TimeEnd = newTimeEnd
	ts.MinuteInfo.Quantity = mb.Seconds
	mb.Seconds = 0

	return
}